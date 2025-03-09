package fips

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
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
	REGISTRATION_NUMBER int    `json:"registration_number"`
	IMAGE_URL           string `json:"image_url"`
	URI                 string `json:"url"`
	AUTHOR              string `json:"author"`
	MAIL                string `json:"mail"`
	DATE                string `json:"registration_date"`
}

type Request struct {
	Search string `json:"search" form:"search"`
	From   string `json:"from" form:"from"`
	To     string `json:"to" form:"to"`
	Limit  int    `json:"limit" form:"limit"`
	Offset int    `json:"offset" form:"offset"`
}

type FipsPage struct {
	url    string
	logger *slog.Logger
}

func New(logger *slog.Logger) *FipsPage {
	return &FipsPage{
		url:    os.Getenv("FIPS_URL"),
		logger: logger,
	}
}

func (a *FipsPage) GetList(c *gin.Context) {
	a.logger.Info("FIPS request", "request", c.Request.URL.Query())
	var r Request
	r.From = c.Request.FormValue("from")
	r.To = c.Request.FormValue("to")
	r.Search = c.Request.FormValue("search")

	strLimit := c.Request.FormValue("limit")
	var err error
	r.Limit, err = strconv.Atoi(strLimit)
	if err != nil {
		r.Limit = -1
	}

	strOffset := c.Request.FormValue("offset")
	r.Offset, err = strconv.Atoi(strOffset)
	if err != nil {
		r.Offset = -1
	}

	a.logger.Info("FIPS request", "request", r)

	byteReq, _ := json.Marshal(r)
	reader := bytes.NewReader(byteReq)

	response, err := http.Post(a.url, "application/json", reader)

	if err != nil {
		a.logger.Error("FIPS request", "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		a.logger.Error("FIPS request", "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var fips FIPS
	err = json.Unmarshal(data, &fips)
	if err != nil {
		a.logger.Error("FIPS request", "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var toWeb []ToWeb
	for _, fip := range fips.Fips {
		toWeb = append(toWeb, ToWeb{
			REGISTRATION_NUMBER: fip.REGISTRATION_NUMBER,
			IMAGE_URL:           fip.IMAGE_URL,
			URI:                 fip.URI,
			AUTHOR:              fip.AUTHOR,
			MAIL:                fip.MAIL,
			DATE:                fip.DATE.Format("02.01.2006"),
		})
	}
	c.JSON(200, gin.H{"fips": toWeb})
}

func (a *FipsPage) GetPage(c *gin.Context) {
	style := "/web/styles/fips-style.css"
	script := "/web/scripts/fips-script.js"
	c.HTML(200, "fips.html", gin.H{
		"title":  "Регистрация товарных знаков",
		"style":  style,
		"script": script,
	})
}
