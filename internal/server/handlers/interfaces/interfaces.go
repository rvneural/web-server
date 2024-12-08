package interfaces

import (
	dbModel "WebServer/internal/models/db/model"

	"github.com/gin-gonic/gin"
)

type DBWorker interface {
	RegisterOperation(uniqID string, operation_type string, user_id int) error
	SetResult(uniqID string, data []byte) error
	GetResult(uniqID string) (dbResult dbModel.DBResult, err error)
	GetAllOperations(limit string, operation_type string, operation_id string) (dbOperations []dbModel.DBResult, err error)
	GetOperationID() (id string, err error)
	GetVersion(uniqID string) (version int64, err error)
	GetUserOperations(user_id int, limit int, operation_type string) (dbResult []dbModel.DBResult, err error)

	CheckForRegistered(email string) bool
	Register(email, hashPassword, FirstName, LastName string) string
	CheckForLogin(email, hashPassword string) (status bool, user_id string)

	GetUserByID(id int) (dbModel.DBUser, error)
	GetAllUsers() ([]dbModel.DBUser, error)
}

type NoResultPage interface {
	GetPage(c *gin.Context, id string)
}
