package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"
)

type StatsPage struct {
}

func New() *StatsPage {
	return &StatsPage{}
}

func (p *StatsPage) GetPage(c *gin.Context) {
	c.JSON(http.StatusOK, stats.Report())
}
