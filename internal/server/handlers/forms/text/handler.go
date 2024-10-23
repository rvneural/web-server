package text

import (
	models "WebServer/internal/models/text"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TextProcessingHandler struct {
	DefaultPrompt string
}

func New(prompt string) *TextProcessingHandler {
	return &TextProcessingHandler{
		DefaultPrompt: prompt,
	}
}

func (n *TextProcessingHandler) HandleForm(c *gin.Context) {
	// Получаем текст и промт
	text := c.Request.FormValue("text")
	var prompt string

	if n.DefaultPrompt == "" {
		prompt = c.Request.FormValue("prompt")
	} else {
		prompt = n.DefaultPrompt
	}

	var resonse models.Response
	var request models.Request
	request.Model = "pro"
	request.Text = text
	request.Prompt = prompt

	if prompt == "{{ rewrite }}" {
		request.Temperature = "0"
	}

	// Маршаллим структуру в JSON и отправляем запрос на Main Server
	data, err := json.Marshal(request)

	if err != nil {
		log.Println(err)

		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusBadRequest, resonse)
		return
	}

	// Получаем ответ от Main Server
	httpRequest, err := http.NewRequest("POST", "http://127.0.0.1:8081/", bytes.NewBuffer(data))

	if err != nil {
		log.Println(err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusBadRequest, resonse)
		return
	}

	httpRequest.Header.Set("Authorization", "Bearer GAuhJOHQ4IQ3sJtFxyRO3OZ84ROyeb")
	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Println(err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusNotAcceptable, resonse)
		return
	}
	defer resp.Body.Close()

	// Парсим ответ от Main Serverа в структуру Response
	byteAns, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusInternalServerError, resonse)
		return
	}

	model := models.Response{}
	err = json.Unmarshal(byteAns, &model)

	if err != nil {
		log.Println(err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusInternalServerError, resonse)
		return
	}

	// Отправляем ответ в виде JSON клиенту
	c.JSON(http.StatusOK, model)
}
