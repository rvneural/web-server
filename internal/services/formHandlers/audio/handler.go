package audio

import (
	models "WebServer/internal/models/audio"
	"io"
	"log"
	"net/http"
	"strings"
)

type RecognitionHandler struct {
}

func New() *RecognitionHandler {
	return &RecognitionHandler{}
}

func (n *RecognitionHandler) handleFileRecognition(r *http.Request) (models.Request, error) {
	var Request models.Request
	file, _, err := r.FormFile("file")  // Полученный файл
	lang := r.FormValue("language")     // Полученный язык
	fileType := r.FormValue("fileType") // Полученный тип файла

	if err != nil {
		log.Println(err)
		return Request, err
	}
	defer file.Close()

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

func (n *RecognitionHandler) handleURLRecognition(r *http.Request) models.Request {
	var Request models.Request
	Request.URL = r.FormValue("url")
	lang := r.FormValue("language")
	Request.Languages = []string{lang}

	return Request
}

// [ ] Audio Handler
func (n *RecognitionHandler) HandleForm(w http.ResponseWriter, r *http.Request) {
	log.Println("Get recognize models from", r.RemoteAddr)
	var Request models.Request
	var err error
	url := r.FormValue("url")
	if url != "" {
		Request = n.handleURLRecognition(r)
	} else {
		Request, err = n.handleFileRecognition(r)
	}

	if err != nil {
		log.Println("Error sending recognition response", err)
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Invalid file", http.StatusBadRequest)
	}

	isDialog := r.FormValue("dialog") // Полученный флаг диалога
	var dialog = false

	if strings.ToLower(isDialog) == "true" {
		dialog = true
	}
	Request.Dialog = dialog

	// Отправляем запрос на распознавание текста
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(n.recognize(Request))

	if err != nil {
		log.Println("Error sending recognition response", err)
	}
}
