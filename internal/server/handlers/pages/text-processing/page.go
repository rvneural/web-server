package text_processing

import (
	"WebServer/internal/services/authorization"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TextProcessingPage struct {
	base string
}

func New() *TextProcessingPage {
	return &TextProcessingPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [x] Text Processing Page
func (rp *TextProcessingPage) GetPage(c *gin.Context) {
	if !authorization.Authorize(c) {
		return
	}

	title := "Обработка текста"
	style := "/web/styles/text-processing-style.js"
	script := "/web/scripts/text-processing-script.js"

	c.HTML(http.StatusOK, "text-processing-page.html", gin.H{
		"title":  title,
		"style":  style,
		"script": script,
	})
}
