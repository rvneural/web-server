package text_rewriting

import (
	"WebServer/internal/services/authorization"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type TextRewritingPage struct {
	base string
}

func New() *TextRewritingPage {
	return &TextRewritingPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

func (rp *TextRewritingPage) GetPage(w http.ResponseWriter, r *http.Request) {
	if !authorization.Authorize(w, r) {
		return
	}
	log.Println("Connection to TextRewritingPage from:", r.RemoteAddr)

	htmlStyle := "../../web/static/pages/text-rewriting/css/style.html"
	content := "../../web/static/pages/text-rewriting/html/page.html"
	script := "../../web/static/pages/text-rewriting/js/script.html"

	t, err := template.ParseFiles(rp.base, htmlStyle, content, script)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}
