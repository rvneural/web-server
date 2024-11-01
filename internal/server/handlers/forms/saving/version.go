package saving

import (
	"WebServer/internal/server/handlers/interfaces"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionSystem struct {
	logger   *slog.Logger
	dbWorher interfaces.DBWorker
}

func NewVersionSystem(dbWorker interfaces.DBWorker, logger *slog.Logger) *VersionSystem {
	return &VersionSystem{
		logger:   logger,
		dbWorher: dbWorker,
	}
}

func (v *VersionSystem) HandleForm(c *gin.Context) {
	operation_id := c.Request.FormValue("id")
	if operation_id == "" {
		v.logger.Error("GetVersion Invalid ID", "id", operation_id)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("incorrect operation ID"),
		})
		return
	}

	version, err := v.dbWorher.GetVersion(operation_id)
	if err != nil {
		v.logger.Error("GetVersion", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	v.logger.Info("GetVersion", "version", version)
	c.JSON(http.StatusOK, gin.H{
		"version": version,
	})

}
