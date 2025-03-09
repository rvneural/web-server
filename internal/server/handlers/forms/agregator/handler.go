package agregator

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsList struct {
	Items []struct {
		Id          uint64 `json:"id"`
		Date        string `json:"date"`
		Title       string `json:"title"`
		Description string `json:"description"`
		IsRt        bool   `json:"isRT"`
	} `json:"items"`
}

type News struct {
	Id          uint64 `json:"id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Rewrite     string `json:"rewrite"`
	Sources     []struct {
		Title       string `json:"title"`
		PubDate     string `json:"pubDate"`
		Name        string `json:"name"`
		Link        string `json:"link"`
		Description string `json:"description"`
		FullText    string `json:"fullText"`
		Enclosure   string `json:"enclosure,omitempty"`
	} `json:"sources"`
}

type Agregator struct {
	url    string
	logger *slog.Logger
}

func New(logger *slog.Logger) *Agregator {
	return &Agregator{
		url:    os.Getenv("AGREGATOR_URL"),
		logger: logger,
	}
}

func (a *Agregator) GetNewsList(c *gin.Context) {
	limit := c.DefaultQuery("limit", "50")
	date := c.DefaultQuery("date", time.Now().Add(3*time.Hour).Format(time.RFC3339))
	q := c.DefaultQuery("q", "")
	url := a.url + "/get?limit=" + limit + "&date=" + date
	a.logger.Info("Q from request", "q", q)
	if q != "" {
		q = strings.ReplaceAll(q, " ", "+")
		url += "&q=" + q
	}
	resp, err := http.Get(url)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		a.logger.Error("Error getting news list", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.logger.Error("Error getting news list", "status", resp.StatusCode)
		c.AbortWithStatus(resp.StatusCode)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		a.logger.Error("Error reading news list", "error", err)
		return
	}
	var newsList NewsList
	err = json.Unmarshal(data, &newsList)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		a.logger.Error("Error unmarshalling news list", "error", err)
		return
	}
	c.JSON(http.StatusOK, newsList)
}

func (a *Agregator) GetNews(c *gin.Context) {
	id := c.Param("id")
	news, err := a.GetNewsData(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, news)
}

func (a *Agregator) GetNewsData(id string) (News, error) {
	url := a.url + "/get/" + id
	resp, err := http.Get(url)
	if err != nil {
		a.logger.Error("Error getting news", "error", err)
		return News{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return News{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Error("Error reading news", "error", err)
		return News{}, err
	}
	var news News
	err = json.Unmarshal(data, &news)
	return news, err
}
