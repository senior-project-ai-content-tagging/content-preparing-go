package scraper

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type SanookScraper struct {
	BaseScraper
}

func (s SanookScraper) Scrap(url string) (string, string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", err
	}

	var content string
	var title string
	doc.Find("h1 .title").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
		log.Println("title: " + title)
	})
	doc.Find("div #EntryReader_0").Each(func(i int, s *goquery.Selection) {
		paragraph := s.Find("p").Text()
		content += paragraph
	})

	return content, title, nil
}

func NewSanookScraper() *SanookScraper {
	return &SanookScraper{
		BaseScraper{Host: "sanook.com"},
	}
}
