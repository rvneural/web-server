package audio

import (
	models "WebServer/internal/models/audio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const addr string = "http://127.0.0.1:8082/"
const key string = "M1PMKzexi0mvX8w7Q1uUz9eH0i3Enw"

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

	httpRequest, err := http.NewRequest("POST", addr, reader)
	if err != nil {
		ans := models.Response{NormText: err.Error(), RawText: err.Error(), Error: err.Error()}
		return ans
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+key)

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
