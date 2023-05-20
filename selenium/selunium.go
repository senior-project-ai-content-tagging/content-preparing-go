package selenium

import (
	"log"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func NewSelenium() (*selenium.WebDriver, error) {
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver("/usr/local/bin/chromedriver"), // Path to chromedriver executable
		selenium.Output(nil), // Disable logging
	}

	// Start Selenium WebDriver service
	_, err := selenium.NewChromeDriverService("/usr/local/bin/chromedriver", 9515, opts...)
	if err != nil {
		log.Fatal("Failed to start Selenium service:", err)
	}

	// Create Chrome capabilities
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless", // Run Chrome in headless mode
			"--no-sandbox",
			"--disable-dev-shm-usage",
		},
	}
	caps.AddChrome(chromeCaps)

	// Start a WebDriver session
	wd, err := selenium.NewRemote(caps, "http://localhost:9515/wd/hub")
	if err != nil {
		log.Fatal("Failed to open session:", err)
	}

	return &wd, err
}
