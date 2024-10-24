package notfound

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotFoundPage struct {
}

func New() *NotFoundPage {
	return &NotFoundPage{}
}

func (n *NotFoundPage) GetPage(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}
