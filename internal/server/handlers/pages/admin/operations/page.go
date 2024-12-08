package operations

import (
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	modelAudio "WebServer/internal/models/db/results/audio"
	modelImage "WebServer/internal/models/db/results/image"
	modelText "WebServer/internal/models/db/results/text"
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
	DATA         []byte `json:"-"`
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
	limit := c.DefaultQuery("limit", "100")
	operation_type := c.DefaultQuery("type", "")
	operation_id := c.DefaultQuery("operation", "")

	operations, err := a.dbWorker.GetAllOperations(limit, operation_type, operation_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       err.Error(),
			"description": "Ошибка при чтении",
		})
	}

	id, err := c.Cookie("user_id")
	if err != nil || id == "" {
		id = "-1"
	}
	var user_id = -1
	if id != "-1" {
		user_id, err = strconv.Atoi(id)
		if err != nil {
			user_id = -1
		}
	}
	var user_status = -1
	if user_id != -1 {
		current_user, err := a.dbWorker.GetUserByID(user_id)
		if err != nil {
			user_status = -1
		} else {
			user_status = current_user.USER_STATUS
		}
	}

	JSONoperations := make([]OperationListElement, len(operations))
	for id, operation := range operations {
		if user_status != -1 && operation.USER_STATUS > user_status {
			continue
		}
		JSONoperations[id] = OperationListElement{
			ID:           operation.ID,
			OPERATION_ID: operation.OPERATION_ID,
			URI:          "/protected/operation/" + operation.OPERATION_ID,
			URL:          "https://" + c.Request.Host + "/protected/operation/" + operation.OPERATION_ID,
			FINISHED:     !operation.IN_PROGRESS,
			TYPE:         operation.OPERATION_TYPE,
			CREATED_AT:   operation.CREATION_DATE.Format("02.01.2006 15:04:05"),
			FINISH_DATE:  operation.FINISH_DATE.Format("02.01.2006 15:04:05"),
			DURATION:     operation.FINISH_DATE.Sub(operation.CREATION_DATE).String(),
			VERSION:      operation.VERSION,
			USER_ID:      operation.USER_ID,
			FIRST_NAME:   operation.FIRST_NAME,
			LAST_NAME:    operation.LAST_NAME,
			EMAIL:        operation.EMAIL,
			DATA:         operation.DATA,
		}
	}

	if strings.Contains(c.GetHeader("Accept"), "text/html") {

		str_max_desc := os.Getenv("MAX_DESC")
		max_desc, err := strconv.Atoi(str_max_desc)
		if err != nil {
			max_desc = 500
		}

		type Operation struct {
			URI         string
			Type        string
			ID          int64
			TITLE       string
			DESCRIPTION string
			AUTHOR      string
			DATE        string
			IMAGE       string
		}

		operations := make([]Operation, 0, len(JSONoperations))
		for _, operation := range JSONoperations {
			var q_operation Operation
			q_operation.URI = operation.URI
			q_operation.Type = operation.TYPE
			q_operation.ID = operation.ID
			q_operation.AUTHOR = operation.FIRST_NAME + " " + operation.LAST_NAME
			q_operation.DATE = operation.CREATED_AT

			switch operation.TYPE {
			case "text":
				data := modelText.DBResult{}
				err := json.Unmarshal(operation.DATA, &data)
				if err != nil {
					q_operation.TITLE = "Ошибка при декодировании"
					q_operation.DESCRIPTION = err.Error()
				} else {
					q_operation.TITLE = data.Prompt
					q_operation.DESCRIPTION = data.OldText
					if len(q_operation.DESCRIPTION) > max_desc {
						q_operation.DESCRIPTION = q_operation.DESCRIPTION[:max_desc]
					}
				}
			case "image":
				data := modelImage.DBResult{}
				err := json.Unmarshal(operation.DATA, &data)
				if err != nil {
					q_operation.TITLE = "Ошибка при декодировании"
					q_operation.DESCRIPTION = err.Error()
				} else {
					q_operation.IMAGE = data.B64string
					q_operation.TITLE = data.Prompt
					q_operation.DESCRIPTION = ""
				}
			case "audio":
				data := modelAudio.DBResult{}
				err := json.Unmarshal(operation.DATA, &data)
				if err != nil {
					q_operation.TITLE = "Ошибка при декодировании"
					q_operation.DESCRIPTION = err.Error()
				} else {
					q_operation.TITLE = "Расшифровано: " + data.FileName
					q_operation.DESCRIPTION = data.NormText
					if len(q_operation.DESCRIPTION) > max_desc {
						q_operation.DESCRIPTION = q_operation.DESCRIPTION[:max_desc]
					}
				}
			}

			operations = append(operations, q_operation)
		}

		c.HTML(http.StatusOK, "admin-operations.html",
			gin.H{
				"style":      "/web/styles/admin-operations.css",
				"title":      "Операции пользователей",
				"Operations": operations,
			})
	} else {
		c.JSON(http.StatusOK, gin.H{"operations": JSONoperations})
	}
}
