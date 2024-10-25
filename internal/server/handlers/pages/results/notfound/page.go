package notfound

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotFoundOperaionPage struct {
}

func New() *NotFoundOperaionPage {
	return &NotFoundOperaionPage{}
}

func (n *NotFoundOperaionPage) GetPage(c *gin.Context, id string) {
	c.HTML(http.StatusNotFound, "no-operation.html", gin.H{
		"title": "Операция не найдена",
		"style": "/web/styles/results/no-operation-style.css",
	})
}
