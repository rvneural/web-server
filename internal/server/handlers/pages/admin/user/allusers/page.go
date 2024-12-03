package allusers

import (
	"WebServer/internal/server/handlers/interfaces"

	"github.com/gin-gonic/gin"
)

type AllUsers struct {
	worker interfaces.DBWorker
}

func New(worker interfaces.DBWorker) *AllUsers {
	return &AllUsers{
		worker: worker,
	}
}

func (u *AllUsers) GetPage(c *gin.Context) {
	users, err := u.worker.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}
