package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthPageHandler struct {
}

func New() *AuthPageHandler {
	return &AuthPageHandler{}
}

func (a *AuthPageHandler) CheckForLogin(c *gin.Context) bool {
	token_str, err := c.Cookie("NeuronNexusAuth")
	if err != nil || len(token_str) == 0 {
		log.Println("None token string")
		return false
	}

	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		log.Println("Invalid token")
		return false
	}
	return true
}

func (a *AuthPageHandler) GetRegisterPage(c *gin.Context) {
	if a.CheckForLogin(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Redirect(http.StatusFound, "/")
		return
	}
	script := "/web/scripts/auth/register.js"
	c.HTML(200, "registration.html", gin.H{
		"title":  "Cозданные изображения",
		"script": script,
	})
}

func (a *AuthPageHandler) GetLoginPage(c *gin.Context) {
	if a.CheckForLogin(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Redirect(http.StatusFound, "/")
		return
	}
	script := "/web/scripts/auth/login.js"
	c.HTML(200, "join.html", gin.H{
		"title":  "Cозданные изображения",
		"script": script,
	})
}

func (a *AuthPageHandler) GetLogoutPage(c *gin.Context) {
	c.SetCookie("NeuronNexusAuth", "", -1, "/", "", false, true)
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Redirect(http.StatusOK, "/")
}
