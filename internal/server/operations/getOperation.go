package operations

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Operation struct {
	lastID map[string]int
}

func New() *Operation {
	lastID := make(map[string]int)

	lastID["image"] = 0
	lastID["audio"] = 0
	lastID["text"] = 0

	return &Operation{lastID: lastID}
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

	id := o.newID(operationType)

	c.JSON(http.StatusOK, gin.H{
		"id":  id,
		"uri": strings.ReplaceAll(c.Request.RequestURI, "/get", "") + "/" + strconv.Itoa(id),
	})
}
