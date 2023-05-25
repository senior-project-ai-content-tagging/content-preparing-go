package scraper

import (
	"fmt"
	"log"
)

type ScraperSelector struct {
	sanookScraper   *SanookScraper
	twitterScraper  *TwitterScraper
	facebookScraper *FacebookScraper
	pantipScraper   *PantipScraper
	siamzoneScraper *SiamzoneScraper
}

func NewSelector(sanookScraper *SanookScraper, twitterScraper *TwitterScraper, facebookScraper *FacebookScraper,
	pantipScraper *PantipScraper,
	siamzoneScraper *SiamzoneScraper,
) *ScraperSelector {
	return &ScraperSelector{
		sanookScraper:   sanookScraper,
		twitterScraper:  twitterScraper,
		facebookScraper: facebookScraper,
		pantipScraper:   pantipScraper,
		siamzoneScraper: siamzoneScraper,
	}
}

func (s ScraperSelector) SelectScraper(inputUrl string) (Scraper, error) {
	listScraper := []Scraper{s.sanookScraper, s.twitterScraper, s.facebookScraper, s.pantipScraper, s.siamzoneScraper}
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
