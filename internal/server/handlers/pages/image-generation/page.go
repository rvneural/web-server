package image_generation

import (
	"WebServer/internal/services/authorization"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ImageGenerationPage struct {
	base string
}

func New() *ImageGenerationPage {
	return &ImageGenerationPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [ ] Image Generation Page
func (rp *ImageGenerationPage) GetPage(w http.ResponseWriter, r *http.Request) {
	if !authorization.Authorize(w, r) {
		return
	}
	log.Println("Connection to ImageGenerationPage from:", r.RemoteAddr)

	htmlStyle := "../../web/static/pages/image-generation/css/style.html"
	content := "../../web/static/pages/image-generation/html/page.html"
	script := "../../web/static/pages/image-generation/js/script.html"

	t, err := template.ParseFiles(rp.base, htmlStyle, content, script)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}
