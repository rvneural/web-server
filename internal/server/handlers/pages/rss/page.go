package rss

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

	response, err := http.Get(r.feedURL)
	if err != nil {
		r.logger.Error("Error while sending request", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	byteResponse, err := io.ReadAll(response.Body)
	if err != nil {
		r.logger.Error("Error while reading response", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = xml.Unmarshal(byteResponse, &rss)
	if err != nil {
		r.logger.Error("Error while unmarshalling response", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.XML(http.StatusOK, rss)

}
