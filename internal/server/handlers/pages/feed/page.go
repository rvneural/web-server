package feed

import (
	"WebServer/internal/models/feed"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type News struct {
	Title       string
	Description string
	URL         string
	Date        string
	Source      string
}

type Page struct {
	logger  *slog.Logger
	feedURL string
}

func New(logger *slog.Logger) *Page {
	RSSURL := os.Getenv("RSS_URL")
	return &Page{
		logger:  logger,
		feedURL: RSSURL,
	}
}

func (r *Page) GetPage(c *gin.Context) {

	var rss feed.RSS

	style := "/web/styles/feed-style.css"
	script := "/web/scripts/feed-script.js"

	response, err := http.Get(r.feedURL)
	if err != nil {
		r.logger.Error("Error while sending request", "error", err)
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"title":  "Актуальные новости (не удалось загрузить)",
			"style":  style,
			"script": script,
		})
		return
	}
	defer response.Body.Close()

	byteResponse, err := io.ReadAll(response.Body)
	if err != nil {
		r.logger.Error("Error while reading response", "error", err)
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"title":  "Актуальные новости (не удалось загрузить)",
			"style":  style,
			"script": script,
		})
		return
	}

	err = xml.Unmarshal(byteResponse, &rss)
	if err != nil {
		r.logger.Error("Error while unmarshalling response", "error", err)
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"title": "Актуальные новости (не удалось загрузить)",
			"style": style,
			"scipr": script,
		})
		return
	}

	news := make([]News, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		news = append(news, News{
			Title:       item.Title,
			Description: item.Description,
			URL:         item.Link,
			Date:        item.PubDate,
			Source:      item.Source,
		})
	}

	c.HTML(http.StatusOK, "feed.html", gin.H{
		"title":  "Актуальные новости",
		"style":  style,
		"script": script,
		"news":   news,
	})

}
