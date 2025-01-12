package scraper

import (
	"time"

	"github.com/gocolly/colly/v2"
)

func NewCollector(maxConcurrency int, delay time.Duration) *colly.Collector {
	domain := "pathofexile\\.com"

	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("www.pathofexile.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainRegexp: domain,
		Parallelism:  maxConcurrency,
		Delay:        delay,
	})

	return c
}
