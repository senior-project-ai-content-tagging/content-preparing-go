package scraper_test

import (
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func testTwitterScraper(t *testing.T) {
	html := `

    `

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	var content string
	articleBox := doc.Find("div[data-testid='tweetText]").First()
	articleBox.Find("span").Each(func(i int, s *goquery.Selection) {
		content += s.Text()
	})

	log.Print(content)
}
