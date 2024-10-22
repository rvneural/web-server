package text_processing

import (
	"WebServer/internal/services/authorization"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Template struct {
	Value string
	Name  string
}

type TextProcessingPage struct {
	base string
}

func New() *TextProcessingPage {
	return &TextProcessingPage{
		base: "../../web/static/pages/template/html/base.html",
	}
}

// [x] Text Processing Page
func (rp *TextProcessingPage) GetPage(c *gin.Context) {
	if !authorization.Authorize(c) {
		return
	}

	title := "Обработка текста"
	style := "/web/styles/text-processing-style.css"
	script := "/web/scripts/text-processing-script.js"

	templates := []Template{
		{
			Value: "0",
			Name:  "Выбрать шаблон",
		},
		{
			Value: "digest",
			Name:  "Дайджест",
		},
		{
			Value: "rewrite",
			Name:  "Рерайт",
		},
		{
			Value: "title",
			Name:  "Заголовок",
		},
		{
			Value: "short",
			Name:  "Сокращение",
		},
		{
			Value: "summary",
			Name:  "Суммаризация",
		},
		{
			Value: "analysis",
			Name:  "Анализ",
		},
		{
			Value: "normalize",
			Name:  "Знаки препинания",
		},
	}

	c.HTML(http.StatusOK, "text-processing-page.html", gin.H{
		"title":     title,
		"style":     style,
		"script":    script,
		"templates": templates,
	})
}
