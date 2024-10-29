package text

import (
	"WebServer/internal/server/handlers/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecognitionResult struct {
	notFoundOperation interfaces.NoResultPage
	progressOperation interfaces.NoResultPage
	dbWorker          interfaces.DBWorker
}

func New(notFoundOperation, progressOperation interfaces.NoResultPage, dbWorker interfaces.DBWorker) *RecognitionResult {
	return &RecognitionResult{
		notFoundOperation: notFoundOperation,
		progressOperation: progressOperation,
		dbWorker:          dbWorker,
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
