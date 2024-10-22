package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	UserName string
	Password string
}

func New(userName, password string) *Auth {
	return &Auth{
		UserName: userName,
		Password: password,
	}
}

// [x] Auth Handler
func (n *Auth) HandleForm(c *gin.Context) {
	login := c.Request.FormValue("login")
	password := c.Request.FormValue("password")

	if login == n.UserName && password == n.Password {
		cookie := &http.Cookie{
			Name:    "login",
			Value:   login,
			Expires: time.Now().Add(time.Hour * 24 * 7), // Cookie expires in 7 days
		}
		http.SetCookie(c.Writer, cookie)

		cookie2 := &http.Cookie{
			Name:    "password",
			Value:   password,
			Expires: time.Now().Add(time.Hour * 24 * 7), // Cookie expires in 7 days
		}
		http.SetCookie(c.Writer, cookie2)

		http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
	} else {
		http.Redirect(c.Writer, c.Request, "/auth", http.StatusSeeOther)
	}
}
