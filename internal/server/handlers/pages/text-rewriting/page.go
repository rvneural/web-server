package text_rewriting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TextRewritingPage struct {
}

func New() *TextRewritingPage {
	return &TextRewritingPage{}
}

func (rp *TextRewritingPage) GetPage(c *gin.Context) {
	title := "Рерайт текста"
	script := "/web/scripts/text-rewriting-script.js"
	style := "/web/styles/text-rewriting-style.css"

	c.HTML(http.StatusOK, "text-rewriting-page.html", gin.H{
		"title":  title,
		"script": script,
		"style":  style,
	})
}
