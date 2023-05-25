package scraper

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

type SiamzoneScraper struct {
	BaseScraper
	wd selenium.WebDriver
}

func (sc SiamzoneScraper) Scrap(url string) (string, string, error) {
	log.Printf("url: %s", url)
	err := sc.wd.Get(url)
	if err != nil {
		return "", "", err
	}

	time.Sleep(5 * time.Second)

	html, err := sc.wd.PageSource()
	if err != nil {
		return "", "", err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", err
	}

	var content string
	var title string
	doc.Find("h1.title").Each(func(i int, s *goquery.Selection) {
		title = strings.TrimSpace(s.Text())
	})
	doc.Find("div.insert_ads").Remove()
	doc.Find("div.ats-slot").Remove()
	doc.Find("script").Remove()
	doc.Find("div.is-size-5-desktop").Each(func(i int, s *goquery.Selection) {
		content += sc.CleanText(s.Text())
	})

	return content, title, nil
}

func NewSiamzoneScraper(webDriver *selenium.WebDriver) *SiamzoneScraper {
	return &SiamzoneScraper{
		BaseScraper: BaseScraper{
			Host: "siamzone.com",
		},
		wd: *webDriver,
	}
}
