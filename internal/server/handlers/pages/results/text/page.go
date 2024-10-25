package text

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
	style := "/web/styles/results/text-processing-style.css"

	id := c.Param("id")

	if len(id) < 10 {
		c.HTML(http.StatusNotFound, "no-operation.html", gin.H{
			"title": "Операция не найдела",
			"style": "/web/styles/results/no-operation-style.css",
		})
		return
	} else if len(id) > 35 {
		c.HTML(http.StatusProcessing, "progress-operation.html", gin.H{
			"title": "Операция еще выполняется",
			"style": "/web/styles/results/progress-operation.css",
		})
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
