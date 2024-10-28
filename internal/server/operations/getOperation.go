package operations

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IDGeneratoer interface {
	Generate() string
}

type Operation struct {
	generator IDGeneratoer
}

func New(generator IDGeneratoer) *Operation {
	return &Operation{generator: generator}
}

// TODO: incorrect name of function
func (o *Operation) GetPage(c *gin.Context) {

	operationType := strings.ToLower(c.Param("type"))

	if operationType != "image" && operationType != "audio" && operationType != "text" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid operation type",
		})
		return
	}

	id := o.generator.Generate()

	/*var urn string
	if strings.ToLower(os.Getenv("TLS_MODE")) == "true" {
		urn = "https://" + config.DOMAIN
	} else {
		urn = "localhost"
	}*/

	c.JSON(http.StatusOK, gin.H{
		"id":  id,
		"url": "/operation/" + operationType + "/" + id,
	})
}
