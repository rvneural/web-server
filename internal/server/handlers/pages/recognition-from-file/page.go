package recognition_from_file

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type RecognitionFromFilePage struct {
	base string
}

func New() *RecognitionFromFilePage {
	return &RecognitionFromFilePage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

func (rp *RecognitionFromFilePage) GetPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to RecognitionFromFilePage from:", r.RemoteAddr)

	htmlStyle := "../../web/static/pages/recognition-from-file/css/style.html"
	content := "../../web/static/pages/recognition-from-file/html/page.html"
	script := "../../web/static/pages/recognition-from-file/js/script.html"

	t, err := template.ParseFiles(rp.base, htmlStyle, content, script)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}
