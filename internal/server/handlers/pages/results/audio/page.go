package audio

import (
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

	raw_text := "Some raw text for ID: " + id
	norm_text := "Some normalized text for ID: " + id

	c.HTML(200, "recognition-result.html", gin.H{
		"title":     "Результаты расшифровки",
		"style":     style,
		"raw_text":  raw_text,
		"norm_text": norm_text,
	})
}
