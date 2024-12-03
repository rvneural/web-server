package operations

import (
	"WebServer/internal/server/handlers/interfaces"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OperationListElement struct {
	ID           int64  `json:"id"`
	OPERATION_ID string `json:"operation_id"`
	URI          string `json:"uri"`
	URL          string `json:"url"`
	FINISHED     bool   `json:"finished"`
	TYPE         string `json:"type"`
	CREATED_AT   string `json:"creation_date"`
	FINISH_DATE  string `json:"finish_date"`
	DURATION     string `json:"duration"`
	VERSION      int64  `json:"version"`
	USER_ID      int    `json:"user_id"`
	FIRST_NAME   string `json:"first_name"`
	LAST_NAME    string `json:"last_name"`
	EMAIL        string `json:"email"`
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

	a.getListOfOperations(c)

}

func (a *AdminOperationListStruct) getListOfOperations(c *gin.Context) {
	limit := c.DefaultQuery("limit", "0")
	operation_type := c.DefaultQuery("type", "")
	operation_id := c.DefaultQuery("operation", "")

	operations, err := a.dbWorker.GetAllOperations(limit, operation_type, operation_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       err.Error(),
			"description": "Ошибка при чтении",
		})
	}

	JSONoperations := make([]OperationListElement, len(operations))
	for id, operation := range operations {
		JSONoperations[id] = OperationListElement{
			ID:           operation.ID,
			OPERATION_ID: operation.OPERATION_ID,
			URI:          "/operation/" + operation.OPERATION_ID,
			URL:          "https://" + c.Request.Host + "/operation/" + operation.OPERATION_ID,
			FINISHED:     !operation.IN_PROGRESS,
			TYPE:         operation.OPERATION_TYPE,
			CREATED_AT:   operation.CREATION_DATE.Add(3 * time.Hour).Format("02.01.2006 15:04:05"),
			FINISH_DATE:  operation.FINISH_DATE.Add(3 * time.Hour).Format("02.01.2006 15:04:05"),
			DURATION:     operation.FINISH_DATE.Sub(operation.CREATION_DATE).String(),
			VERSION:      operation.VERSION,
			USER_ID:      operation.USER_ID,
			FIRST_NAME:   operation.FIRST_NAME,
			LAST_NAME:    operation.LAST_NAME,
			EMAIL:        operation.EMAIL,
		}
	}

	c.JSON(http.StatusOK, gin.H{"operations": JSONoperations})
}
