package audio

import (
	config "WebServer/internal/config/services/audio2text-service"
	models "WebServer/internal/models/audio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func (n *RecognitionHandler) recognize(request models.Request) models.Response {
	// Превращаем структуру в JSON-строку
	data, err := json.Marshal(request)

	// В случае ошибки — возвращаем ошибку на морду
	if err != nil {
		ans := models.Response{NormText: err.Error(), RawText: err.Error(), Error: err.Error()}
		return ans
	}

	// Создаем reader и отправляем POST запрос на сервер
	reader := bytes.NewReader(data)

	httpRequest, err := http.NewRequest("POST", config.URL, reader)
	if err != nil {
		ans := models.Response{NormText: err.Error(), RawText: err.Error(), Error: err.Error()}
		return ans
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resonse, err := client.Do(httpRequest)

	// В случае ошибки — возвращаем ошибку на морду
	if err != nil {
		ans := models.Response{NormText: err.Error(), RawText: err.Error(), Error: err.Error()}
		return ans
	}
	defer resonse.Body.Close()

	// Читаем из BODY запроса массив байт
	ans, err := io.ReadAll(resonse.Body)

	if err != nil {
		ans := models.Response{NormText: err.Error(), RawText: err.Error(), Error: err.Error()}
		return ans
	}

	s_ans := models.Response{}
	err = json.Unmarshal(ans, &s_ans)
	if err != nil {
		ans := models.Response{NormText: err.Error(), RawText: err.Error(), Error: err.Error()}
		return ans
	}
	return s_ans
}
