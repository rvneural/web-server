package audio

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NoResultPage interface {
	GetPage(c *gin.Context, id string)
}

type RecognitionResult struct {
	notFoundOperation NoResultPage
	progressOperation NoResultPage
}

func New(notFoundOperation, progressOperation NoResultPage) *RecognitionResult {
	return &RecognitionResult{
		notFoundOperation: notFoundOperation,
		progressOperation: progressOperation,
	}
}

func (r *RecognitionResult) GetPage(c *gin.Context) {

	style := "/web/styles/results/recognition-style.css"

	id := c.Param("id")

	if len(id) < 10 {
		r.notFoundOperation.GetPage(c, id)
		return
	} else if len(id) > 35 {
		r.progressOperation.GetPage(c, id)
		return
	}

	raw_text := "Some raw text for ID: " + id
	norm_text := "Some normalized text for ID: " + id

	c.HTML(http.StatusOK, "recognition-result.html", gin.H{
		"title":     "Результаты расшифровки",
		"style":     style,
		"raw_text":  raw_text,
		"norm_text": norm_text,
	})
}
