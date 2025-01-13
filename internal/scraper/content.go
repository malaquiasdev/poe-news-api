package scraper

import (
	"poe-news-api/internal/database"
	"poe-news-api/internal/util"
	"strings"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

func ScrapeContent(c *colly.Collector, posts *[]database.Post) (*[]database.Post, error) {
	c.OnRequest(func(r *colly.Request) {
		util.Logger.Info("Visiting URL",
			zap.String("url", r.URL.String()),
		)
	})

	c.OnError(func(r *colly.Response, err error) {
		util.Logger.Error("Request error",
			zap.String("url", r.Request.URL.String()),
			zap.Error(err),
		)
	})

	c.OnHTML("tr.staff", func(e *colly.HTMLElement) {
		rawMessage := e.ChildText("div.content")
		rawMessage = strings.ReplaceAll(rawMessage, `\n`, "\n")
		paragraphs := strings.Split(rawMessage, "\n\n")
		(*posts)[0].FirstMessage = e.Text
		util.Logger.Info("Scraping content", zap.Any("content", paragraphs))
	})

	for post := 0; post < len(*posts); post++ {
		util.Logger.Info("Queueing post visit", zap.String("url", (*posts)[post].Link))
		c.Visit((*posts)[post].Link)
	}

	c.Wait()

	util.Logger.Info("Scraping completed successfully")
	return posts, nil
}
