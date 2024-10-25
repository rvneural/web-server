package audio

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

	style := "/web/styles/results/recognition-style.css"

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

	raw_text := "Some raw text for ID: " + id
	norm_text := "Some normalized text for ID: " + id

	c.HTML(http.StatusOK, "recognition-result.html", gin.H{
		"title":     "Результаты расшифровки",
		"style":     style,
		"raw_text":  raw_text,
		"norm_text": norm_text,
	})
}
