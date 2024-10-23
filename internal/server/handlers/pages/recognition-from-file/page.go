package recognition_from_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Language struct {
	Code string
	Name string
}

type RecognitionFromFilePage struct {
}

func New() *RecognitionFromFilePage {
	return &RecognitionFromFilePage{}
}

func (rp *RecognitionFromFilePage) GetPage(c *gin.Context) {
	style := "/web/styles/recognition-style.css"
	script := "/web/scripts/recognition-script.js"
	title := "Расшифровка аудио и видео"

	laguages := []Language{
		{
			Code: "ru-RU",
			Name: "Русский",
		},
		{
			Code: "auto",
			Name: "Определить автоматически",
		},
		{
			Code: "en-US",
			Name: "Английский",
		},
		{
			Code: "de-DE",
			Name: "Немецкий",
		},
		{
			Code: "es-ES",
			Name: "Испанский",
		},
		{
			Code: "fi-FI",
			Name: "Финский",
		},
		{
			Code: "fr-FR",
			Name: "Французский",
		},
		{
			Code: "pl-PL",
			Name: "Польский",
		},
		{
			Code: "he-HE",
			Name: "Иврит",
		},
		{
			Code: "it-IT",
			Name: "Итальянский",
		},
		{
			Code: "kk-KZ",
			Name: "Казахский",
		},
		{
			Code: "nl-NL",
			Name: "Голландский",
		},
		{
			Code: "pt-PT",
			Name: "Португальский",
		},
		{
			Code: "pt-BR",
			Name: "Бразильский португальский",
		},
		{
			Code: "sv-SE",
			Name: "Шведский",
		},
		{
			Code: "tr-TR",
			Name: "Турецкий",
		},
		{
			Code: "uz-UZ",
			Name: "Узбецкий (латиница)",
		},
	}

	c.HTML(http.StatusOK, "recognition-page.html", gin.H{
		"title":     title,
		"style":     style,
		"script":    script,
		"languages": laguages,
	})
}
