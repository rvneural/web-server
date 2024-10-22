package image_generation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageGenerationPage struct {
	base string
}

type Rations struct {
	Value string
	Name  string
}

func New() *ImageGenerationPage {
	return &ImageGenerationPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [x] Image Generation Page
func (rp *ImageGenerationPage) GetPage(c *gin.Context) {
	style := "/web/styles/image-generation-style.css"
	script := "/web/scripts/image-generation-script.js"
	title := "Генерация изображений"

	rations := []Rations{
		{
			Value: "3-2",
			Name:  "3:2",
		},
		{
			Value: "1-1",
			Name:  "1:1",
		},
		{
			Value: "16-9",
			Name:  "16:9",
		},
		{
			Value: "9-16",
			Name:  "9:16",
		},
		{
			Value: "2-3",
			Name:  "2:3",
		},
	}

	c.HTML(http.StatusOK, "image-generation-page.html", gin.H{
		"title":   title,
		"style":   style,
		"script":  script,
		"rations": rations,
	})
}
