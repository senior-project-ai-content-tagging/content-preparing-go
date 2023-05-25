package scraper

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

type FacebookScraper struct {
	BaseScraper
	wd selenium.WebDriver
}

func (s FacebookScraper) Scrap(url string) (string, string, error) {
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
	doc.Find("div[data-testid='post_message']").Each(func(i int, s *goquery.Selection) {
		content += s.Text()
	})

	return content, title, nil
}

func NewFacebookScraper(webDriver *selenium.WebDriver) *FacebookScraper {
	return &FacebookScraper{
		BaseScraper: BaseScraper{
			Host: "facebook.com",
		},
		wd: *webDriver,
	}
}
