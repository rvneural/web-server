package saving

import (
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"log/slog"
	"net/http"

	audioResult "WebServer/internal/models/db/results/audio"
	textResult "WebServer/internal/models/db/results/text"

	"github.com/gin-gonic/gin"
)

type SavingSystem struct {
	dbWorker interfaces.DBWorker
	logger   *slog.Logger
}

func New(dbWorker interfaces.DBWorker, logger *slog.Logger) *SavingSystem {
	return &SavingSystem{
		dbWorker: dbWorker,
		logger:   logger,
	}
}

func (s *SavingSystem) HandleForm(c *gin.Context) {
	operation_id := c.Request.FormValue("id")
	operation_type := c.Request.FormValue("type")

	s.logger.Info("Запрс на сохранение", "params", c.Request.Form)

	if operation_id == "" || operation_type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
	}
	switch operation_type {
	case "audio":
		result := audioResult.DBResult{
			FileName: c.Request.FormValue("file_name"),
			RawText:  c.Request.FormValue("raw_text"),
			NormText: c.Request.FormValue("norm_text"),
		}
		byteResult, err := json.Marshal(result)
		if err != nil {
			s.logger.Error("Error while saving to db", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		err = s.dbWorker.SetResult(operation_id, byteResult)
		if err != nil {
			s.logger.Error("Error while saving to db", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	case "text":
		result := textResult.DBResult{
			OldText: c.Request.FormValue("old_text"),
			NewText: c.Request.FormValue("new_text"),
			Prompt:  c.Request.FormValue("prompt"),
		}
		byteResult, err := json.Marshal(result)
		if err != nil {
			s.logger.Error("Error while saving to db", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		err = s.dbWorker.SetResult(operation_id, byteResult)
		if err != nil {
			s.logger.Error("Error while saving to db", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
