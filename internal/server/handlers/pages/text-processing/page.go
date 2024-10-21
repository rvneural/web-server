package text_processing

import (
	"WebServer/internal/services/authorization"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type TextProcessingPage struct {
	base string
}

func New() *TextProcessingPage {
	return &TextProcessingPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [ ] Text Processing Page
func (rp *TextProcessingPage) GetPage(w http.ResponseWriter, r *http.Request) {
	if !authorization.Authorize(w, r) {
		return
	}
	log.Println("Connection to TextProcessingPage from:", r.RemoteAddr)

	htmlStyle := "../../web/static/pages/text-processing/css/style.html"
	content := "../../web/static/pages/text-processing/html/page.html"
	script := "../../web/static/pages/text-processing/js/script.html"

	t, err := template.ParseFiles(rp.base, htmlStyle, content, script)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}
