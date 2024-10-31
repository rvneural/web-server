package photopea

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PhotopeaHandler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *PhotopeaHandler {
	return &PhotopeaHandler{
		logger: logger,
	}
}

func (p *PhotopeaHandler) HandleForm(c *gin.Context) {
	p.logger.Info("Redirection to photopea", "user", c.RemoteIP())
	c.Redirect(http.StatusSeeOther, "https://www.photopea.com/")
}
