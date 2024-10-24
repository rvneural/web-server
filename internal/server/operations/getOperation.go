package operations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Operation struct {
}

func New() *Operation {
	return &Operation{}
}

func (o *Operation) GetNewID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id": o.newID(),
	})
}
