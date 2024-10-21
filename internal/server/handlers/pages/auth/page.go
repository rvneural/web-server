package auth

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Auth struct {
}

func New() *Auth {
	return &Auth{}
}

// [ ] Auth page
func (rp *Auth) GetPage(w http.ResponseWriter, r *http.Request) {

	log.Println("Connection to Auth from:", r.RemoteAddr)

	page := "../../web/static/pages/auth/auth.html"

	t, err := template.ParseFiles(page)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.Execute(w, nil)
}
