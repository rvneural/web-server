package text

import (
	config "WebServer/internal/config/services/text2text-service"
	models "WebServer/internal/models/text"
	"log/slog"
	"strconv"
	"strings"

	"WebServer/internal/server/handlers/interfaces"

	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TextProcessingHandler struct {
	DefaultPrompt string
	dbWorker      interfaces.DBWorker
	logger        *slog.Logger
}

func New(prompt string, dbWorker interfaces.DBWorker, logger *slog.Logger) *TextProcessingHandler {
	return &TextProcessingHandler{
		DefaultPrompt: prompt,
		dbWorker:      dbWorker,
		logger:        logger,
	}
}

func (n *TextProcessingHandler) HandleForm(c *gin.Context) {

	id := c.Request.FormValue("id")
	id = strings.TrimSpace(id)
	str_id, err := c.Cookie("user_id")
	var user_id int
	if err != nil {
		n.logger.Error("Getting user id", "error", err)
		user_id = 0
	} else {
		user_id, err = strconv.Atoi(str_id)
		if err != nil {
			n.logger.Error("Converting user id", "error", err)
			user_id = 0
		}
	}

	n.logger.Info("Operation TEXT from", "user", user_id, "str", str_id)

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
	request.OperationId = id
	request.UserID = user_id

	if prompt == "{{ rewrite }}" {
		request.Temperature = "0.3"
	}

	// Маршаллим структуру в JSON и отправляем запрос на Main Server
	data, err := json.Marshal(request)

	if err != nil {
		n.logger.Error("Marshalling request", "error", err)

		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusBadRequest, resonse)
		return
	}

	// Получаем ответ от Main Server
	httpRequest, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(data))

	if err != nil {
		n.logger.Error("Creating request", "error", err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusBadRequest, resonse)
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		n.logger.Error("Sending request", "error", err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusBadGateway, resonse)
		return
	}
	defer resp.Body.Close()

	// Парсим ответ от Main Serverа в структуру Response
	byteAns, err := io.ReadAll(resp.Body)

	if err != nil {
		n.logger.Error("Reading response", "error", err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusInternalServerError, resonse)
		return
	}

	model := models.Response{}
	err = json.Unmarshal(byteAns, &model)

	if err != nil {
		n.logger.Error("Unmarshalling response", "error", err)
		resonse.OldText = text
		resonse.NewText = err.Error()
		c.JSON(http.StatusInternalServerError, resonse)
		return
	}

	// Отправляем ответ в виде JSON клиенту
	c.JSON(http.StatusOK, model)
}
