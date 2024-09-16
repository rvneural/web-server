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

func (n *RecognitionHandler) HandleForm(w http.ResponseWriter, r *http.Request) {
	log.Println("Get recognize models from", r.RemoteAddr)
	var Request models.Request

	file, _, err := r.FormFile("file")  // Полученный файл
	lang := r.FormValue("language")     // Полученный язык
	fileType := r.FormValue("fileType") // Полученный тип файла
	isDialog := r.FormValue("dialog")   // Полученный флаг диалога
	var dialog = false

	if strings.ToLower(isDialog) == "true" {
		dialog = true
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)

	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	// В зависимости от того, прошло сжатие данных успешно или нет
	// выполняется передача данных на Main Server в сжатом
	// или не в сжатом виде
	Request.FileData = fileData
	Request.Dialog = dialog
	Request.Languages = []string{lang}
	Request.FileType = fileType

	// Отправляем запрос на распознавание текста
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(n.recognize(Request))

	if err != nil {
		log.Println("Error sending recognition response", err)
	}
}
