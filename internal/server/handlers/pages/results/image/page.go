package image

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
