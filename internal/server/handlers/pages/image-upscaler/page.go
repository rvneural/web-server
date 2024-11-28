package imageupscaler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageUpscalerPage struct {
}

func New() *ImageUpscalerPage {
	return &ImageUpscalerPage{}
}

func (rp *ImageUpscalerPage) GetPage(c *gin.Context) {
	style := "/web/styles/image-upscale-style.css"
	script := "/web/scripts/image-upscale-script.js"
	title := "Обработка изображения"

	c.HTML(http.StatusOK, "image-upscale-page.html", gin.H{
		"title":  title,
		"style":  style,
		"script": script,
	})
}
