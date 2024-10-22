package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
}

func New() *Auth {
	return &Auth{}
}

// [x] Auth page
func (rp *Auth) GetPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth-page.html", gin.H{
		"title": "Авторизация",
	})
}
