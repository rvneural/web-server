package image

import "github.com/gin-gonic/gin"

type RecognitionResult struct {
}

func New() *RecognitionResult {
	return &RecognitionResult{}
}

func (r *RecognitionResult) GetPage(c *gin.Context) {
	c.HTML(200, "image-generation-result.html", gin.H{
		"title": "Main website",
	})
}
