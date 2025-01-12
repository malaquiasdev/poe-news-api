package scraper

import (
	"fmt"
	"poe-news-api/internal/database"
	"poe-news-api/internal/util"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

func ScrapeForum(c *colly.Collector, startURL string, pageLimit int) ([]database.Post, error) {
	var posts []database.Post

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

	c.OnHTML("tr", func(e *colly.HTMLElement) {

		post := database.Post{
			Title:      e.ChildText("td.thread > div.thread_title"),
			Author:     e.ChildText("td.thread > div.postBy > span.post_by_account > a"),
			PostDate:   e.ChildText("td.thread > div.postBy > span.post_date"),
			Link:       e.Request.AbsoluteURL(e.ChildAttr("td.thread > div.thread_title > div.title > a", "href")),
			NumReplies: e.Request.AbsoluteURL(e.ChildAttr("td.views > div > span", "href")),
		}

		if post.Title != "" && post.Link != "" {
			posts = append(posts, post)

		}

		util.Logger.Info("Scraped",
			zap.Any("post", post),
		)
	})

	for page := 1; page <= pageLimit; page++ {
		targetURL := fmt.Sprintf("%s?page=%d", startURL, page)
		util.Logger.Info("Queueing page visit", zap.String("url", targetURL))

		c.Visit(targetURL)
	}

	c.Wait()

	util.Logger.Info("Scraping completed successfully", zap.Int("num_posts", len(posts)))
	return posts, nil
}
