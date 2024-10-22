package image_generation

import (
	"WebServer/internal/services/authorization"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageGenerationPage struct {
	base string
}

func New() *ImageGenerationPage {
	return &ImageGenerationPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [x] Image Generation Page
func (rp *ImageGenerationPage) GetPage(c *gin.Context) {
	if !authorization.Authorize(c) {
		return
	}

	style := "/web/styles/image-generation-style.css"
	script := "/web/scripts/image-generation-script.js"
	title := "Генерация изображений"

	c.HTML(http.StatusOK, "image-generation-page.html", gin.H{
		"title":  title,
		"style":  style,
		"script": script,
	})
}
