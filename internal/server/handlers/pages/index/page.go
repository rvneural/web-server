package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Page struct {
}

func New() *Page {
	return &Page{}
}

func (p *Page) GetPage(c *gin.Context) {
	style := "/web/styles/index.css"
	script := "/web/scripts/index.js"
	c.HTML(http.StatusOK, "index.html", gin.H{
		"style":  style,
		"script": script,
		"title":  "Главная страница",
	})
}
