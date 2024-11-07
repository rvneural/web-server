package operations

import (
	"WebServer/internal/server/handlers/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IDGeneratoer interface {
	Generate() string
}

type Operation struct {
	dbWorker interfaces.DBWorker
}

func New(dbWorker interfaces.DBWorker) *Operation {
	return &Operation{
		dbWorker: dbWorker,
	}
}

func (o *Operation) GetPage(c *gin.Context) {
	id, err := o.dbWorker.GetOperationID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":  id,
		"url": "/operation/" + id,
	})
}
