package text

import (
	config "WebServer/internal/config/services/text2text-service"
	models "WebServer/internal/models/text"
	"strings"

	dbModel "WebServer/internal/models/db/results/text"
	"WebServer/internal/server/handlers/interfaces"

	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TextProcessingHandler struct {
	DefaultPrompt string
	dbWorker      interfaces.DBWorker
}

func New(prompt string, dbWorker interfaces.DBWorker) *TextProcessingHandler {
	return &TextProcessingHandler{
		DefaultPrompt: prompt,
		dbWorker:      dbWorker,
	}
}

func (n *TextProcessingHandler) HandleForm(c *gin.Context) {

	id := c.Request.FormValue("id")
	id = strings.TrimSpace(id)
	var dbError error
	if len(id) != 0 {
		dbError = n.dbWorker.RegisterOperation(id, "text")
	}

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
		request.Temperature = "0.3"
	}

	// Маршаллим структуру в JSON и отправляем запрос на Main Server
	data, err := json.Marshal(request)

	if err != nil {
		log.Println(err)

		resonse.OldText = text
		resonse.NewText = err.Error()
		go n.saveErrorToDB(id, err.Error(), prompt, text)
		c.JSON(http.StatusBadRequest, resonse)
		return
	}

	// Получаем ответ от Main Server
	httpRequest, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(data))

	if err != nil {
		log.Println(err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		go n.saveErrorToDB(id, err.Error(), prompt, text)
		c.JSON(http.StatusBadRequest, resonse)
		return
	}

	httpRequest.Header.Set("Authorization", "Bearer "+config.KEY)
	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Println(err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		go n.saveErrorToDB(id, err.Error(), prompt, text)
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
		go n.saveErrorToDB(id, err.Error(), prompt, text)
		c.JSON(http.StatusInternalServerError, resonse)
		return
	}

	model := models.Response{}
	err = json.Unmarshal(byteAns, &model)

	if err != nil {
		log.Println(err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		go n.saveErrorToDB(id, err.Error(), prompt, text)
		c.JSON(http.StatusInternalServerError, resonse)
		return
	}

	go func(id string, dbError error) {
		if dbError == nil && len(id) != 0 {
			ddbResult := dbModel.DBResult{
				OldText: model.OldText,
				NewText: model.NewText,
				Prompt:  prompt,
			}
			byteDbResult, err := json.Marshal(ddbResult)
			if err == nil {
				n.dbWorker.SetResult(id, byteDbResult)
			}
		}
	}(id, dbError)

	// Отправляем ответ в виде JSON клиенту
	c.JSON(http.StatusOK, model)
}

func (n *TextProcessingHandler) saveErrorToDB(id string, errorMsg string, prompt string, oldText string) {
	if len(id) == 0 {
		return
	}
	dbResult := dbModel.DBResult{
		OldText: oldText,
		NewText: errorMsg,
		Prompt:  prompt,
	}
	byteDbResult, err := json.Marshal(dbResult)
	if err == nil {
		n.dbWorker.SetResult(id, byteDbResult)
	}
}
