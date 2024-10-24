package operations

import "github.com/gin-gonic/gin"

type Operation struct {
}

func New() *Operation {
	return &Operation{}
}

func (o *Operation) GetNewID(c *gin.Context) {

}
