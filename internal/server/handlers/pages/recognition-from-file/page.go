package recognition_from_file

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type RecognitionFromFilePage struct {
}

func New() *RecognitionFromFilePage {
	return &RecognitionFromFilePage{}
}

func (rp *RecognitionFromFilePage) GetPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to homePage from:", r.RemoteAddr)
	t, err := template.ParseFiles(
		"../../../../web/static/pages/recognition-from-file/html/page.html",
		"../../../../web/static/pages/template/html/base.html",
		"pages/styles/style.html",
		"pages/scripts/style.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}
