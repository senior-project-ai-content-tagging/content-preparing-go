package scraper

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

type PantipScraper struct {
	BaseScraper
	wd selenium.WebDriver
}

func (sc PantipScraper) Scrap(url string) (string, string, error) {
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
	// write html to file
	f, err := os.Create("facebook.html")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	f.WriteString(html)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", err
	}

	var content string
	var title string
	doc.Find(".display-post-title").Each(func(i int, s *goquery.Selection) {
		title = strings.TrimSpace(s.Text())
	})
	doc.Find("div .display-post-story").Each(func(i int, s *goquery.Selection) {
		content += sc.CleanText(s.Text())
	})

	return content, title, nil
}

func NewPantipScraper(webDriver *selenium.WebDriver) *PantipScraper {
	return &PantipScraper{
		BaseScraper: BaseScraper{
			Host: "pantip.com",
		},
		wd: *webDriver,
	}
}
