package image

import (
	"net/http"
	"strconv"

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
	int_id, err := strconv.Atoi(id)
	if err != nil || int_id < 1 {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}

	prompt := "Some prompt for ID: " + id
	seed := 321321321321
	image := "/web/static/img/templates/1-1.png"

	c.HTML(200, "image-generation-result.html", gin.H{
		"title":  "Результаты генерации",
		"style":  style,
		"prompt": prompt,
		"seed":   seed,
		"image":  image,
	})
}
