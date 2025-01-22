package fips

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type FIPS struct {
	Fips []struct {
		REGISTRATION_NUMBER int       `json:"registration_number"`
		IMAGE_URL           string    `json:"image_url"`
		URI                 string    `json:"url"`
		AUTHOR              string    `json:"author"`
		MAIL                string    `json:"mail"`
		DATE                time.Time `json:"registration_date"`
	} `json:"fips"`
}

type ToWeb struct {
	REGISTRATION_NUMBER int
	IMAGE_URL           string
	URI                 string
	AUTHOR              string
	MAIL                string
	DATE                string
}

type FipsPage struct {
	url string
}

func New() *FipsPage {
	return &FipsPage{
		url: os.Getenv("FIPS_URL"),
	}
}

func (a *FipsPage) GetPage(c *gin.Context) {
	search := c.Query("search")
	if search != "" {
		search = "?search=" + search
	}

	from := c.DefaultQuery("from", "")
	if from != "" {
		if search != "" {
			from = "&from=" + from
		} else {
			from = "?from=" + from
		}
	}

	to := c.DefaultQuery("to", "")
	if to != "" {
		if from != "" || search != "" {
			to = "&to=" + to
		} else {
			to = "?to=" + to
		}
	}

	response, err := http.Get(a.url + search + from + to)

	if err != nil {
		style := "/web/styles/fips-style.css"
		script := "/web/scripts/fips-script.js"
		c.HTML(200, "fips.html", gin.H{
			"title":  "Регистрация товарных знаков",
			"style":  style,
			"script": script,
		})
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		style := "/web/styles/fips-style.css"
		script := "/web/scripts/fips-script.js"
		c.HTML(200, "fips.html", gin.H{
			"title":  "Регистрация товарных знаков",
			"style":  style,
			"script": script,
		})
		return
	}

	var fips FIPS
	err = json.Unmarshal(data, &fips)
	if err != nil {
		style := "/web/styles/fips-style.css"
		script := "/web/scripts/fips-script.js"
		c.HTML(200, "fips.html", gin.H{
			"title":  "Регистрация товарных знаков",
			"style":  style,
			"script": script,
		})
		return
	}

	content := fips.Fips
	var toWeb []ToWeb
	for _, fip := range content {
		toWeb = append(toWeb, ToWeb{
			REGISTRATION_NUMBER: fip.REGISTRATION_NUMBER,
			IMAGE_URL:           fip.IMAGE_URL,
			URI:                 fip.URI,
			AUTHOR:              fip.AUTHOR,
			MAIL:                fip.MAIL,
			DATE:                fip.DATE.Format("02.01.2006"),
		})
	}
	style := "/web/styles/fips-style.css"
	script := "/web/scripts/fips-script.js"
	c.HTML(200, "fips.html", gin.H{
		"title":  "Регистрация товарных знаков",
		"Fips":   toWeb,
		"style":  style,
		"script": script,
	})
}
