package interfaces

import (
	dbModel "WebServer/internal/models/db/model"

	"github.com/gin-gonic/gin"
)

type DBWorker interface {
	RegisterOperation(uniqID string) error
	SetResult(uniqID string, data []byte) error
	GetResult(uniqID string) (dbResult dbModel.DBResult, err error)
}

type NoResultPage interface {
	GetPage(c *gin.Context, id string)
}