package image

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecognitionResult struct {
}

func New() *RecognitionResult {
	return &RecognitionResult{}
}

func (r *RecognitionResult) GetPage(c *gin.Context) {

	style := "/web/styles/results/image-generation-style.css"

	id := c.Param("id")

	if len(id) < 10 {
		c.HTML(http.StatusNotFound, "no-operation.html", gin.H{
			"title": "Операция не найдена",
			"style": "/web/styles/results/no-operation-style.css",
		})
		return
	} else if len(id) > 35 {
		c.HTML(http.StatusNotFound, "progress-operation.html", gin.H{
			"title": "Операция еще выполняется",
			"style": "/web/styles/results/progress-operation.css",
		})
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
