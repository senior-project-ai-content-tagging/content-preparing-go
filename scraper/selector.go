package scraper

import (
	"fmt"
	"log"
)

type ScraperSelector struct {
	sanookScraper  *SanookScraper
	twitterScraper *TwitterScraper
}

func NewSelector(sanookScraper *SanookScraper, twitterScraper *TwitterScraper) *ScraperSelector {
	return &ScraperSelector{
		sanookScraper:  sanookScraper,
		twitterScraper: twitterScraper,
	}
}

func (s ScraperSelector) SelectScraper(inputUrl string) (Scraper, error) {
	listScraper := []Scraper{s.sanookScraper, s.twitterScraper}
	for _, scraper := range listScraper {
		validScraper, err := scraper.CheckDomain(inputUrl)
		if err != nil {
			return nil, err
		}

		if validScraper {
			log.Printf("scraper: %v", scraper)
			return scraper, nil
		}
	}

	return nil, fmt.Errorf("Not allow domain %s", inputUrl)
}
