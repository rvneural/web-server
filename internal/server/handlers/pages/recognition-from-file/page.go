package recognition_from_file

import (
	"WebServer/internal/services/authorization"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecognitionFromFilePage struct {
	base string
}

func New() *RecognitionFromFilePage {
	return &RecognitionFromFilePage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [x] Recognition Page
func (rp *RecognitionFromFilePage) GetPage(c *gin.Context) {
	if !authorization.Authorize(c) {
		return
	}

	style := "/web/styles/recognition-style.css"
	script := "/web/scripts/recognition-script.js"
	title := "Расшифровка аудио и видео"

	c.HTML(http.StatusOK, "recognition-page.html", gin.H{
		"title":  title,
		"style":  style,
		"script": script,
	})
}
