package mediafeed

import (
	"WebServer/internal/models/feed"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

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
	logger         *slog.Logger
	feedURL        string
	mu             sync.Mutex
	LastNews       []News
	MaxDescription int
}

func New(logger *slog.Logger) *Page {
	RSSURL := os.Getenv("MEDIA_RSS_URL")

	max_description_str := os.Getenv("MAX_DESC")
	max_description, err := strconv.Atoi(max_description_str)
	if err != nil {
		max_description = 255
	}

	page := &Page{
		logger:         logger,
		feedURL:        RSSURL,
		LastNews:       nil,
		mu:             sync.Mutex{},
		MaxDescription: max_description,
	}

	go func() {
		for true {
			page.UpdateNews()
			time.Sleep(1 * time.Minute)
		}
	}()

	return page
}

func (p *Page) UpdateNews() error {
	var rss feed.RSS
	response, err := http.Get(p.feedURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	byteResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(byteResponse, &rss)
	if err != nil {
		return err
	}
	var news []News

	for _, item := range rss.Channel.Items {

		if len(item.Description) > p.MaxDescription {
			item.Description = item.Description[:p.MaxDescription] + "..."
		}

		date, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			item.PubDate = date.Format("02.01.2006 15:04")
		}

		news = append(news, News{
			Title:       item.Title,
			Description: item.Description,
			URL:         item.Link,
			Date:        item.PubDate,
			Source:      item.Source,
		})
	}

	p.mu.Lock()
	p.LastNews = news
	p.mu.Unlock()
	return nil
}

func (r *Page) GetPage(c *gin.Context) {

	style := "/web/styles/feed-style.css"
	script := "/web/scripts/media-script.js"

	if r.LastNews == nil {
		err := r.UpdateNews()
		if err != nil {
			r.logger.Error("Error while reading response", "error", err)
			c.HTML(http.StatusOK, "feed.html", gin.H{
				"title":  "Новости СМИ (не удалось загрузить)",
				"style":  style,
				"script": script,
			})
			return
		}
	}

	r.mu.Lock()
	news := r.LastNews
	r.mu.Unlock()

	c.HTML(http.StatusOK, "feed.html", gin.H{
		"title":  "Новости СМИ",
		"style":  style,
		"script": script,
		"news":   news,
	})

}
