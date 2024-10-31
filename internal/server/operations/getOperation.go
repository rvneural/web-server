package operations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IDGeneratoer interface {
	Generate() string
}

type Operation struct {
	generator IDGeneratoer
}

func New(generator IDGeneratoer) *Operation {
	return &Operation{
		generator: generator,
	}
}

func (o *Operation) GetPage(c *gin.Context) {
	id := o.generator.Generate()

	c.JSON(http.StatusOK, gin.H{
		"id":  id,
		"url": "/operation/" + id,
	})
}
