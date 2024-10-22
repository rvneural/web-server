package imageupscaler

import (
	"WebServer/internal/services/authorization"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageUpscalerPage struct {
	base string
}

func New() *ImageUpscalerPage {
	return &ImageUpscalerPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [x] Image Upscale Page
func (rp *ImageUpscalerPage) GetPage(c *gin.Context) {
	if !authorization.Authorize(c) {
		return
	}

	style := "/web/styles/image-upscale-style.css"
	script := "/web/scripts/image-upscale-script.js"
	title := "Увеличение изображения"

	c.HTML(http.StatusOK, "image-upscale-page.html", gin.H{
		"title":  title,
		"style":  style,
		"script": script,
	})
}
