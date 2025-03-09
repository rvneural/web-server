package text_processing

import (
	"encoding/json"
	"net/http"

	config "WebServer/internal/config/services/text2text-service"

	"github.com/gin-gonic/gin"
)

type Template struct {
	Value string
	Name  string
}

type TemplateModel struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type TemplatesModel struct {
	Templates []TemplateModel `json:"templates"`
}

type TextProcessingPage struct {
}

func New() *TextProcessingPage {
	return &TextProcessingPage{}
}

func (rp *TextProcessingPage) GetPage(c *gin.Context) {
	title := "Обработка текста"
	style := "/web/styles/text-processing-style.css"
	script := "/web/scripts/text-processing-script.js"

	var tmpls TemplatesModel

	req, err := http.NewRequest("GET", config.URL, nil)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&tmpls)
		}
	}

	templates := []Template{
		{
			Value: "0",
			Name:  "Выбрать шаблон",
		},
	}

	for _, tmpls := range tmpls.Templates {
		templates = append(templates, Template{
			Value: tmpls.Name,
			Name:  tmpls.Title,
		})
	}

	c.HTML(http.StatusOK, "text-processing-page.html", gin.H{
		"title":     title,
		"style":     style,
		"script":    script,
		"templates": templates,
	})
}
