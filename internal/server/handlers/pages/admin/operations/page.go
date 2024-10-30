package operations

import (
	"WebServer/internal/server/handlers/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OperationListElement struct {
	ID           int    `json:"id"`
	OPERATION_ID string `json:"operation_id"`
	URI          string `json:"uri"`
	FINISHED     bool   `json:"finished"`
	TYPE         string `json:"type"`
}

type AllOperations struct {
	Operations []OperationListElement `json:"operations"`
}

type AdminOperationListStruct struct {
	dbWorker interfaces.DBWorker
}

func New(dbWorker interfaces.DBWorker) *AdminOperationListStruct {
	return &AdminOperationListStruct{
		dbWorker: dbWorker,
	}
}

func (a *AdminOperationListStruct) GetPage(c *gin.Context) {
	operations, err := a.dbWorker.GetAllOperations()
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
