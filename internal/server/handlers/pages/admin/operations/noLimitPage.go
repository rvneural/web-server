package operations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NoLimitAllResults struct {
}

func NewNoLimit() *NoLimitAllResults {
	return &NoLimitAllResults{}
}

func (n *NoLimitAllResults) GetPage(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/admin/operations/0")
}
