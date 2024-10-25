package text

import (
	"github.com/gin-gonic/gin"
)

type RecognitionResult struct {
}

func New() *RecognitionResult {
	return &RecognitionResult{}
}

func (r *RecognitionResult) GetPage(c *gin.Context) {
	style := "/web/styles/results/text-processing-style.css"

	id := c.Param("id")

	old_text := "Some old text for ID: " + id
	new_text := "Some new text for ID: " + id
	prompt := "Some prompt for ID: " + id

	c.HTML(200, "text-processing-result.html", gin.H{
		"title":    "Результаты обработки",
		"style":    style,
		"old_text": old_text,
		"new_text": new_text,
		"prompt":   prompt,
	})
}
