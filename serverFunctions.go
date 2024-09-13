package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"transcriptorWeb/scripts"
)

// Функция обработки запроса в распознавание текста
func recognizeFromWeb(w http.ResponseWriter, r *http.Request) {
	// Если метод GET — возвращаем ошибку
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	log.Println("Get recognize request from", r.RemoteAddr)

	file, _, err := r.FormFile("file")  // Полученный файл
	lang := r.FormValue("language")     // Полученный язык
	fileType := r.FormValue("fileType") // Полученный тип файла
	s_dialog := r.FormValue("dialog")   // Полученный флаг диалога
	var dialog bool = false

	if strings.ToLower(s_dialog) == "true" {
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

	// Буффер для сжатого файла
	var compressedFile bytes.Buffer

	log.Println("Compressing file")
	// Сжимаем файл
	writer := zlib.NewWriter(&compressedFile)
	_, err = writer.Write(fileData)
	writer.Close()

	// В зависимости от того, прошло сжатие данных успешно или нет
	// выполняется передача данных на Main Server в сжатом
	// или не в сжатом виде
	var Request scripts.Request
	Request.FileData = fileData
	Request.Dialog = dialog
	Request.Languages = []string{lang}
	Request.FileType = fileType

	w.WriteHeader(http.StatusOK)
	w.Write(scripts.Recognize(Request))
}

// Функция обработки запроса в переписывание текста
func rewriteFromWeb(w http.ResponseWriter, r *http.Request) {
	log.Println("New rewrite from web request from", r.RemoteAddr)

	// Получаем текст для переписывания
	text := r.FormValue("text")

	// Структура для запроса к Main Server
	type RequestForRewrite struct {
		Text   string `json:"text"`
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	// Структура ответа от Main Server
	type Response struct {
		OldText string `json:"oldText"`
		NewText string `json:"newText"`
	}

	var request RequestForRewrite
	request.Text = text
	request.Prompt = "{{ rewrite }}"
	request.Model = "lite"

	// Отправляем запрос к Main Server
	data, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}

	// Получаем ответ от Main Server
	httpRequest, err := http.NewRequest("POST", "http://127.0.0.1:8081/", bytes.NewBuffer(data))

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}

	httpRequest.Header.Set("Authorization", "Bearer GAuhJOHQ4IQ3sJtFxyRO3OZ84ROyeb")
	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}
	defer resp.Body.Close()

	// Парсим ответ от Main Serverа
	byteAns, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}

	// Отправляем полученный ответ клиенту
	w.WriteHeader(http.StatusOK)
	w.Write(byteAns)
}

// Функция обработки запроса по распознаванию текста
func processTextFromWeb(w http.ResponseWriter, r *http.Request) {
	log.Println("New process text from web request from", r.RemoteAddr)

	// Получаем текст и промт
	text := r.FormValue("text")
	prompt := r.FormValue("prompt")

	// Структура для обмена данными с Main Serverом
	type RequestForRewrite struct {
		Text   string `json:"text"`
		Prompt string `json:"prompt"`
	}

	// Структура для ответа от Main Serverа, которая в дальнейшем парсится в JS
	type Response struct {
		OldText string `json:"oldText"`
		NewText string `json:"newText"`
	}

	var request RequestForRewrite
	request.Text = text
	request.Prompt = prompt

	// Маршаллим структуру в JSON и отправляем запрос на Main Server
	data, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}

	// Получаем ответ от Main Server
	httpRequest, err := http.NewRequest("POST", "http://127.0.0.1:8081/", bytes.NewBuffer(data))

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}

	httpRequest.Header.Set("Authorization", "Bearer GAuhJOHQ4IQ3sJtFxyRO3OZ84ROyeb")
	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}
	defer resp.Body.Close()

	// Парсим ответ от Main Serverа в структуру Response
	byteAns, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		var resonse Response
		resonse.OldText = text
		resonse.NewText = err.Error()
		s_data, _ := json.Marshal(resonse)
		w.WriteHeader(200)
		w.Write(s_data)
		return
	}

	// Отправляем ответ в виде JSON клиенту
	w.WriteHeader(http.StatusOK)
	w.Write(byteAns)
}

// Функция обработки запрос на генерацию изображения
func generateImageFromWeb(w http.ResponseWriter, r *http.Request) {
	log.Println("New generate image from web request from", r.RemoteAddr)

	prompt := r.FormValue("prompt")
	seed := r.FormValue("seed")
	widthRatio := r.FormValue("widthRatio")
	heightRatio := r.FormValue("heightRatio")

	type RequestForImage struct {
		Prompt      string `json:"prompt"`
		Seed        string `json:"seed"`
		WidthRatio  string `json:"widthRatio"`
		HeightRatio string `json:"heightRatio"`
	}

	var request RequestForImage
	request.Prompt = prompt
	request.Seed = seed
	request.WidthRatio = widthRatio
	request.HeightRatio = heightRatio

	byteRequets, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	httpRequest, err := http.NewRequest("POST", "http://127.0.0.1:8083/", bytes.NewBuffer(byteRequets))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer KU2WTZBzFWn4Ko9lJ7TlpmUXwkHc8Y")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	defer resp.Body.Close()

	byteResp, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(byteResp)
}
