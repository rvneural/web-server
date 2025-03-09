package text2speech

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
	c.HTML(http.StatusOK, "text2speech.html", nil)
}
