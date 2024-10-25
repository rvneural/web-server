package text

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
	style := "/web/styles/results/text-processing-style.css"

	id := c.Param("id")

	if len(id) < 10 {
		r.notFoundOperation.GetPage(c, id)
		return
	} else if len(id) > 35 {
		r.progressOperation.GetPage(c, id)
		return
	}

	old_text := "Some old text for ID: " + id
	new_text := "Some new text for ID: " + id
	prompt := "Some prompt for ID: " + id

	c.HTML(http.StatusOK, "text-processing-result.html", gin.H{
		"title":    "Результаты обработки",
		"style":    style,
		"old_text": old_text,
		"new_text": new_text,
		"prompt":   prompt,
	})
}
