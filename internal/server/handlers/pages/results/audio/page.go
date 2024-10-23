package audio

import "github.com/gin-gonic/gin"

type RecognitionResult struct {
}

func New() *RecognitionResult {
	return &RecognitionResult{}
}

func (r *RecognitionResult) GetPage(c *gin.Context) {
	c.HTML(200, "recognition-result.html", gin.H{
		"title": "Main website",
	})
}
