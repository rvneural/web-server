package price

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Page struct {
	logger   *slog.Logger
	priceURL string
}

func New(logger *slog.Logger) *Page {
	return &Page{
		logger:   logger,
		priceURL: os.Getenv("PRICE_PARSER_URL"),
	}
}

func (r *Page) GetPage(c *gin.Context) {

	style := "/web/styles/prices.css"
	script := "/web/scripts/prices.js"

	models, err := http.Get(r.priceURL + "/models")
	var umnarshalledModels Models
	if err == nil {
		data, err2 := io.ReadAll(models.Body)
		if err2 == nil {
			json.Unmarshal(data, &umnarshalledModels)
		}
	}

	queryModel := c.DefaultQuery("model", "")
	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

	var table string = ""

	if queryModel != "" {
		var url string = r.priceURL + "/prices/" + queryModel
		if from != "" {
			url += "?from=" + from
			if to != "" {
				url += "&to=" + to
			}
		} else if from == "" && to != "" {
			url += "?to=" + to
		}

		prices, err := http.Get(url)

		if err == nil {
			var Product ParsedProduct
			data, err2 := io.ReadAll(prices.Body)
			if err2 == nil {
				json.Unmarshal(data, &Product)
				table = r.jsonToHTMLTable(&Product)
			} else {
				r.logger.Error(err2.Error())
			}
		} else {
			r.logger.Error(err.Error())
		}
	}

	var fromValue string = ""
	var toValue string = ""

	if from != "" {
		fromDate, err := time.Parse("02.01.2006", from)
		if err == nil {
			fromValue = fromDate.Format("2006-01-02")
		}
	}
	if to != "" {
		toDate, err := time.Parse("02.01.2006", to)
		if err == nil {
			toValue = toDate.Format("2006-01-02")
		}
	}

	c.HTML(http.StatusOK, "prices.html", gin.H{
		"title":  "Парсер цен",
		"style":  style,
		"script": script,
		"Models": umnarshalledModels.Models,
		"Table":  table,
		"Value":  queryModel,
		"From":   fromValue,
		"To":     toValue,
	})

}
