package image

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

	style := "/web/styles/results/image-generation-style.css"

	id := c.Param("id")

	if len(id) < 10 {
		r.notFoundOperation.GetPage(c, id)
		return
	} else if len(id) > 35 {
		r.progressOperation.GetPage(c, id)
		return
	}

	prompt := "Some prompt for ID: " + id
	seed := 321321321321
	image := "/web/static/img/templates/9-16.png"

	c.HTML(http.StatusOK, "image-generation-result.html", gin.H{
		"title":  "Результаты генерации",
		"style":  style,
		"prompt": prompt,
		"seed":   seed,
		"image":  image,
	})
}
