package audio

import (
	models "WebServer/internal/models/audio"
	dbModel "WebServer/internal/models/db/results/audio"
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
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
	var dbError error
	if len(id) != 0 {
		var u_id int
		user_id, err := c.Cookie("user_id")
		if err != nil {
			u_id = 0
		} else {
			u_id, err = strconv.Atoi(user_id)
			if err != nil {
				u_id = 0
			}
		}
		dbError = n.dbWorker.RegisterOperation(id, "audio", u_id)
	}

	var Request models.Request
	var err error = nil
	url := c.Request.FormValue("url")
	if url != "" {
		Request = n.handleURLRecognition(c)
	} else {
		Request, err = n.handleFileRecognition(c)
		if err != nil {
			go n.logger.Error("Error sending recognition response", "error", err)
			go n.saveErrorToDB(id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	Request.Dialog = (strings.ToLower(c.Request.FormValue("dialog")) == "true")

	n.logger.Info("FILE DIALOG", "dialog", Request.Dialog, "form", c.Request.FormValue("dialog"))

	response := n.recognize(Request)

	response.NormText = strings.TrimSpace(response.NormText)
	response.RawText = strings.TrimSpace(response.RawText)

	go func(id string, dbError error) {
		if dbError == nil && len(id) != 0 {
			dbResult := dbModel.DBResult{
				RawText:  strings.TrimSpace(response.RawText),
				NormText: strings.TrimSpace(response.NormText),
				FileName: strings.TrimSpace(filename),
			}
			byteData, err := json.Marshal(dbResult)
			if err == nil {
				n.dbWorker.SetResult(id, byteData)
			}
		}
	}(id, dbError)

	// Отправляем запрос на распознавание текста
	c.JSON(http.StatusOK, response)
}
