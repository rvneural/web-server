package text

import "github.com/gin-gonic/gin"

type RecognitionResult struct {
}

func New() *RecognitionResult {
	return &RecognitionResult{}
}

func (r *RecognitionResult) GetPage(c *gin.Context) {
	style := "/web/styles/results/text-processing-style.css"
	c.HTML(200, "text-processing-result.html", gin.H{
		"title": "Результаты обработки",
		"style": style,
	})
}
