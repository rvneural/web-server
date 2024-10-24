package audio

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

	style := "/web/styles/results/recognition-style.css"

	id := c.Param("id")
	int_id, err := strconv.Atoi(id)
	if err != nil || int_id < 1 {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}

	raw_text := "Some raw text for ID: " + id
	norm_text := "Some normalized text for ID: " + id

	c.HTML(200, "recognition-result.html", gin.H{
		"title":     "Результаты расшифровки",
		"style":     style,
		"raw_text":  raw_text,
		"norm_text": norm_text,
	})
}
