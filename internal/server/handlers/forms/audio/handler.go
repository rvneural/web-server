package audio

import (
	models "WebServer/internal/models/audio"
	"WebServer/internal/server/handlers/interfaces"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RecognitionHandler struct {
	dbWorker interfaces.DBWorker
	logger   *slog.Logger
}

func New(dbWorker interfaces.DBWorker, logger *slog.Logger) *RecognitionHandler {
	return &RecognitionHandler{
		dbWorker: dbWorker,
		logger:   logger,
	}
}

func (n *RecognitionHandler) handleFileRecognition(c *gin.Context) (models.Request, error) {
	var Request models.Request

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return Request, err
	}
	defer file.Close()

	lang := c.Request.FormValue("language")
	fileType := c.Request.FormValue("fileType")
	fileData, err := io.ReadAll(file)
	if err != nil {
		n.logger.Error("Error reading file", "error", err)
		return Request, err
	}

	if c.Request.FormValue("whisper") == "true" {
		Request.Model = "whisper"
	} else {
		Request.Model = "yandex"
	}
	Request.File.Data = fileData
	Request.Languages = []string{lang}
	Request.File.Type = fileType

	return Request, nil
}

func (n *RecognitionHandler) handleURLRecognition(c *gin.Context) models.Request {
	var Request models.Request
	Request.URL = c.Request.FormValue("url")
	lang := c.Request.FormValue("language")
	Request.Languages = []string{lang}

	return Request
}

func (n *RecognitionHandler) HandleForm(c *gin.Context) {
	id := c.Request.FormValue("id")
	id = strings.TrimSpace(id)
	filename := c.Request.FormValue("filename")
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

	n.logger.Info("Operation AUDIO from", "user", user_id, "str", str_id)

	var Request models.Request
	url := c.Request.FormValue("url")
	if url != "" {
		Request = n.handleURLRecognition(c)
	} else {
		Request, err = n.handleFileRecognition(c)
		if err != nil {
			go n.logger.Error("Error sending recognition response", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	Request.Dialog = (strings.ToLower(c.Request.FormValue("dialog")) == "true")
	Request.OperationId = id
	Request.File.Name = filename
	Request.UserID = user_id

	n.logger.Info("FILE DIALOG", "dialog", Request.Dialog, "form", c.Request.FormValue("dialog"))

	response := n.recognize(Request)

	response.NormText = strings.TrimSpace(response.NormText)
	response.RawText = strings.TrimSpace(response.RawText)

	// Отправляем запрос на распознавание текста
	c.JSON(http.StatusOK, response)
}
