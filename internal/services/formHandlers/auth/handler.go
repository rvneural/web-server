package auth

import (
	"net/http"
	"time"
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

// [ ] Auth Handler
func (n *Auth) HandleForm(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == n.UserName && password == n.Password {
		cookie := &http.Cookie{
			Name:    "login",
			Value:   login,
			Expires: time.Now().Add(time.Hour * 24 * 7), // Cookie expires in 7 days
		}
		http.SetCookie(w, cookie)

		cookie2 := &http.Cookie{
			Name:    "password",
			Value:   password,
			Expires: time.Now().Add(time.Hour * 24 * 7), // Cookie expires in 7 days
		}
		http.SetCookie(w, cookie2)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
	}
}
