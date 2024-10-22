package imageupscaler

import (
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

func (rp *ImageUpscalerPage) GetPage(c *gin.Context) {
	style := "/web/styles/image-upscale-style.css"
	script := "/web/scripts/image-upscale-script.js"
	title := "Увеличение изображения"

	c.HTML(http.StatusOK, "image-upscale-page.html", gin.H{
		"title":  title,
		"style":  style,
		"script": script,
	})
}
