package photopea

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PhotopeaHandler struct {
}

func New() *PhotopeaHandler {
	return &PhotopeaHandler{}
}

func (p *PhotopeaHandler) HandleForm(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "https://www.photopea.com/")
}
