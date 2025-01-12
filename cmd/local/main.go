package main

import (
	"flag"
	"fmt"
	"time"

	"go.uber.org/zap"

	"poe-news-api/internal/scraper"
	"poe-news-api/internal/util"
)

func main() {

	// Command-line flags for flexibility in local runs
	startURL := flag.String("url", "https://www.pathofexile.com/forum/view-forum/2212", "Path of Exile forum URL to scrape")
	pageLimit := flag.Int("pages", 1, "Number of pages to scrape")
	concurrency := flag.Int("concurrency", 1, "Max concurrency for the scraper")
	delayMillis := flag.Int("delay", 500, "Delay in milliseconds between requests")

	flag.Parse()

	c := scraper.NewCollector(
		*concurrency,
		time.Duration(*delayMillis)*time.Millisecond,
	)

	util.Logger.Info("Starting forum scrape",
		zap.String("url", *startURL),
		zap.Int("pageLimit", *pageLimit),
		zap.Int("concurrency", *concurrency),
		zap.Int("delayMillis", *delayMillis),
	)

	if err := scraper.ScrapeForum(c, *startURL, *pageLimit); err != nil {
		util.Logger.Fatal("Failed to scrape forum", zap.Error(err))
	}

	util.Logger.Info("Scraping session finished")
	fmt.Println("Scraping complete! Check logs for details.")
}
