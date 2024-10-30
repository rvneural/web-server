package operations

import (
	"WebServer/internal/server/handlers/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LimitOperationsWithType struct {
	dbWorker interfaces.DBWorker
}

func NewLimitOperationsWithType(dbWorker interfaces.DBWorker) *LimitOperationsWithType {
	return &LimitOperationsWithType{dbWorker: dbWorker}
}

func (l *LimitOperationsWithType) GetPage(c *gin.Context) {
	var limit int
	str_limit := c.Param("limit")
	limit, err := strconv.Atoi(str_limit)
	if err != nil || limit < 0 {
		c.Redirect(http.StatusSeeOther, "/admin/operations/0")
		return
	}
	operation_type := c.Param("type")

	operations, err := l.dbWorker.GetAllOperations(limit, operation_type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	JSONoperations := make([]OperationListElement, len(operations))

	for id, operation := range operations {
		JSONoperations[id] = OperationListElement{
			ID:           operation.ID,
			OPERATION_ID: operation.OPERATION_ID,
			URI:          "/operation/" + operation.OPERATION_ID,
			FINISHED:     !operation.IN_PROGRESS,
			TYPE:         operation.OPERATION_TYPE,
		}
	}

	c.JSON(http.StatusOK, AllOperations{
		Operations: JSONoperations,
	})
}
