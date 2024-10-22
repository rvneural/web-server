package audio

import (
	models "WebServer/internal/models/audio"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RecognitionHandler struct {
}

func New() *RecognitionHandler {
	return &RecognitionHandler{}
}

func (n *RecognitionHandler) handleFileRecognition(c *gin.Context) (models.Request, error) {
	var Request models.Request
	file, _, err := c.Request.FormFile("file") // Полученный файл
	if err != nil {
		return Request, err
	}
	defer file.Close()

	lang := c.Request.FormValue("language")     // Полученный язык
	fileType := c.Request.FormValue("fileType") // Полученный тип файла
	fileData, err := io.ReadAll(file)

	if err != nil {
		log.Println(err)
		return Request, err
	}

	// В зависимости от того, прошло сжатие данных успешно или нет
	// выполняется передача данных на Main Server в сжатом
	// или не в сжатом виде
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

// [x] Audio Handler
func (n *RecognitionHandler) HandleForm(c *gin.Context) {
	var Request models.Request
	var err error = nil
	url := c.Request.FormValue("url")
	if url != "" {
		Request = n.handleURLRecognition(c)
	} else {
		Request, err = n.handleFileRecognition(c)
	}

	if err != nil {
		log.Println("Error sending recognition response", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	isDialog := c.Request.FormValue("dialog") // Полученный флаг диалога
	var dialog = false

	if strings.ToLower(isDialog) == "true" {
		dialog = true
	}
	Request.Dialog = dialog

	// Отправляем запрос на распознавание текста
	c.JSON(http.StatusOK, n.recognize(Request))
}
