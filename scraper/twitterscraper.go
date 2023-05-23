package scraper

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

type TwitterScraper struct {
	BaseScraper
	wd selenium.WebDriver
}

func (s TwitterScraper) Scrap(url string) (string, string, error) {
	log.Printf("url: %s", url)
	err := s.wd.Get(url)
	if err != nil {
		return "", "", err
	}

	time.Sleep(5 * time.Second)

	html, err := s.wd.PageSource()
	if err != nil {
		return "", "", err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", err
	}

	var content string
	title := ""
	articleBox := doc.Find("div[data-testid='tweetText']").First()
	articleBox.Find("span").Each(func(i int, s *goquery.Selection) {
		log.Println(s)
		content += s.Text()
	})

	return content, title, nil
}

func NewTwitterScraper(webDriver *selenium.WebDriver) *TwitterScraper {
	return &TwitterScraper{
		BaseScraper: BaseScraper{
			Host: "twitter.com",
		},
		wd: *webDriver,
	}
}
