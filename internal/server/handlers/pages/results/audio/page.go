package audio

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

	style := "/web/styles/results/recognition-style.css"

	id := c.Param("id")

	res, err := r.dbWorker.GetResult(id)

	if err != nil {
		r.notFoundOperation.GetPage(c, id)
		return
	} else if res.IN_PROGRESS {
		r.progressOperation.GetPage(c, id)
		return
	}

	c.HTML(http.StatusOK, "recognition-result.html", gin.H{
		"title":     "Результаты расшифровки",
		"style":     style,
		"raw_text":  raw_text,
		"norm_text": norm_text,
	})
}
