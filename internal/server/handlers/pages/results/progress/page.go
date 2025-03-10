package progress

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgressOperaionPage struct {
}

func New() *ProgressOperaionPage {
	return &ProgressOperaionPage{}
}

func (n *ProgressOperaionPage) GetPage(c *gin.Context, id string) {
	c.HTML(http.StatusNotFound, "progress-operation.html", gin.H{
		"title": "Операция еще выполняется",
		"id":    id,
		"url":   c.Request.URL.String(),
		"style": "/web/styles/results/progress-operation-style.css",
	})
}
