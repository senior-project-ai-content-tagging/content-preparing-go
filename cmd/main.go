package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/senior-project-ai-content-tagging/content-preparing/config"
	"github.com/senior-project-ai-content-tagging/content-preparing/entity"
	"github.com/senior-project-ai-content-tagging/content-preparing/publisher"
	"github.com/senior-project-ai-content-tagging/content-preparing/repository"
	"github.com/senior-project-ai-content-tagging/content-preparing/scraper"
	"github.com/senior-project-ai-content-tagging/content-preparing/selenium"
	"github.com/senior-project-ai-content-tagging/content-preparing/translator"
)

func main() {
	mux := http.NewServeMux()
	handler, err := initHanler()
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/submit-weblink", handler.submitWeblink)
	mux.HandleFunc("/submit-content", handler.submitContent)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type Handler struct {
	db                *sqlx.DB
	contentRepository repository.ContentRepository
	ticketRepository  repository.TicketRepository
	translator        translator.Translator
	sanookScraper     *scraper.SanookScraper
	twitterScraper    *scraper.TwitterScraper
	selector          *scraper.ScraperSelector
}

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func initHanler() (*Handler, error) {
	dbConfig := config.GetDatabaseConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSL)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	azureConfig := config.GetAzureConfig()
	azTranslator := translator.NewAzureTranslator(azureConfig.ApiKey)

	sanookScraper := scraper.NewSanookScraper()

	wd, err := selenium.NewSelenium()
	if err != nil {
		return nil, err
	}

	twitterScraper := scraper.NewTwitterScraper(wd)

	selector := scraper.NewSelector(sanookScraper, twitterScraper)

	contentRepository := repository.NewContentRepository()
	ticketRepository := repository.NewTicketRepository()

	return &Handler{
		db:                db,
		contentRepository: contentRepository,
		ticketRepository:  ticketRepository,
		translator:        azTranslator,
		sanookScraper:     sanookScraper,
		selector:          selector,
	}, nil
}

func (h Handler) submitWeblink(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var m PubSubMessage
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	// byte slice unmarshalling handles base64 decoding.
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	var ticketSubmit entity.TicketSubmitWeblinkPubSub
	err = json.Unmarshal(m.Message.Data, &ticketSubmit)
	if err != nil {
		log.Printf("json.Unmarshal to struct: %v", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	scraper, err := h.selector.SelectScraper(ticketSubmit.Url)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	content, title, err := scraper.Scrap(ticketSubmit.Url)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	contentEN, err := h.translator.Translate(content)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}
	titleEN, err := h.translator.Translate(title)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	tx := h.db.MustBegin()

	sql, entityTicket := h.ticketRepository.UpdateTicketStatus(entity.Ticket{Id: ticketSubmit.TicketId, Status: "PROCESSING"})
	tx.NamedExecContext(ctx, sql, entityTicket)

	sql, ticketId := h.ticketRepository.FindTicketById(ticketSubmit.TicketId)
	var ticket entity.Ticket
	err = tx.GetContext(ctx, &ticket, sql, ticketId)

	newContent := entity.Content{
		Id:        ticket.ContentId,
		TitleTH:   title,
		ContentTH: content,
		TitleEN:   titleEN,
		ContentEN: contentEN,
	}
	log.Println(newContent)

	sql, entityContent := h.contentRepository.UpdateContent(newContent)
	tx.NamedExecContext(ctx, sql, entityContent)

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	pubsubConfig := config.GetPublisherConfig()
	prepapredPublisher, err := publisher.NewPublisher(ctx, pubsubConfig.ProjectID, pubsubConfig.PreparedTopicID)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	err = prepapredPublisher.Publish(entity.TicketPubSub{TicketId: ticketId})
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	res.WriteHeader(http.StatusOK)
}

func (h Handler) submitContent(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var m PubSubMessage
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	// byte slice unmarshalling handles base64 decoding.
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	var ticketSubmit entity.TicketContentPubSub
	err = json.Unmarshal(m.Message.Data, &ticketSubmit)
	if err != nil {
		log.Printf("json.Unmarshal to struct: %v", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	tx := h.db.MustBegin()
	sql, entityTicket := h.ticketRepository.UpdateTicketStatus(entity.Ticket{Id: ticketSubmit.TicketId, Status: "PROCESSING"})
	tx.NamedExecContext(ctx, sql, entityTicket)

	sql, ticketId := h.ticketRepository.FindTicketById(ticketSubmit.TicketId)
	var ticket entity.Ticket
	err = tx.GetContext(ctx, &ticket, sql, ticketId)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	sql, contentId := h.contentRepository.FindContentById(ticket.ContentId)
	var content entity.Content
	err = tx.GetContext(ctx, &content, sql, contentId)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	contentEN, err := h.translator.Translate(content.ContentTH)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	titleEN, err := h.translator.Translate(content.TitleTH)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	newContent := entity.Content{
		Id:        ticket.ContentId,
		TitleTH:   content.TitleTH,
		ContentTH: content.ContentTH,
		TitleEN:   titleEN,
		ContentEN: contentEN,
	}
	log.Println(newContent)

	sql, entityContent := h.contentRepository.UpdateContent(newContent)
	tx.NamedExecContext(ctx, sql, entityContent)

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	pubsubConfig := config.GetPublisherConfig()
	prepapredPublisher, err := publisher.NewPublisher(ctx, pubsubConfig.ProjectID, pubsubConfig.PreparedTopicID)
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	err = prepapredPublisher.Publish(entity.TicketPubSub{TicketId: ticketId})
	if err != nil {
		log.Print(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	res.WriteHeader(http.StatusOK)
}