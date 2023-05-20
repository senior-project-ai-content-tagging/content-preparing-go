package scraper

import (
	"net/url"
	"strings"
)

type Scraper interface {
	Scrap(url string) (string, string, error)
	CheckDomain(inputUrl string) (bool, error)
}

type BaseScraper struct {
	Host string
}

func (s BaseScraper) CheckDomain(inputUrl string) (bool, error) {
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		return false, err
	}

	host := parsedUrl.Host
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}

	return host == s.Host, nil
}
