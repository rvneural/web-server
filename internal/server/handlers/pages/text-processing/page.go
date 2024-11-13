package text_processing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Template struct {
	Value string
	Name  string
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
			Value: "correct",
			Name:  "Исправление ошибок",
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
