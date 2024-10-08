package imageupscaler

import (
	"WebServer/internal/services/authorization"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ImageUpscalerPage struct {
	base string
}

func New() *ImageUpscalerPage {
	return &ImageUpscalerPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

func (rp *ImageUpscalerPage) GetPage(w http.ResponseWriter, r *http.Request) {
	if !authorization.Authorize(w, r) {
		return
	}
	log.Println("Connection to ImageUpscalerPage from:", r.RemoteAddr)

	htmlStyle := "../../web/static/pages/image-upscaler/css/style.html"
	content := "../../web/static/pages/image-upscaler/html/page.html"
	script := "../../web/static/pages/image-upscaler/js/script.html"

	t, err := template.ParseFiles(rp.base, htmlStyle, content, script)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}
