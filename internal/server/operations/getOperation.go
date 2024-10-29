package operations

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type IDGeneratoer interface {
	Generate() string
}

type Operation struct {
	generator     IDGeneratoer
	acceptedTypes []string
}

func New(generator IDGeneratoer) *Operation {
	return &Operation{
		generator:     generator,
		acceptedTypes: []string{"image", "audio", "text", "test"},
	}
}

// TODO: incorrect name of function
func (o *Operation) GetPage(c *gin.Context) {

	operationType := strings.ToLower(c.Param("type"))

	if !slices.Contains(o.acceptedTypes, operationType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid operation type",
		})
		return
	}

	id := o.generator.Generate()

	c.JSON(http.StatusOK, gin.H{
		"id":  id,
		"url": "/operation/" + operationType + "/" + id,
	})
}
