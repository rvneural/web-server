package scripts

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Адрес, по которому доступен Main Server по протоколу HTTP
const addr string = "http://127.0.0.1:8082/"
const key string = "M1PMKzexi0mvX8w7Q1uUz9eH0i3Enw"

// Функция расшифровки файла, приходящего из веб-морды
func Recognize(request Request) []byte {
	log.Println("New request to recognition service")
	// Превращаем структуру в JSON-строку
	data, err := json.Marshal(request)

	// В случае ошибки — возвращаем ошибку на морду
	if err != nil {
		log.Println(err)
		ans := Answer{NormText: err.Error(), RawText: err.Error()}
		byteAns, _ := json.Marshal(ans)
		return byteAns
	}

	// Создаем reader и отправляем POST запрос на сервер
	reader := bytes.NewReader(data)

	httpRequest, err := http.NewRequest("POST", addr, reader)
	if err != nil {
		log.Println(err)
		ans := Answer{NormText: err.Error(), RawText: err.Error()}
		byteAns, _ := json.Marshal(ans)
		return byteAns
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+key)

	client := http.Client{}

	resonse, err := client.Do(httpRequest)

	// В случае ошибки — возвращаем ошибку на морду
	if err != nil {
		log.Println(err)
		ans := Answer{NormText: err.Error(), RawText: err.Error()}
		byteAns, _ := json.Marshal(ans)
		return byteAns
	}
	defer resonse.Body.Close()
	log.Println("Response from recognition service:", resonse)

	// Если статус результата не 200 — обращаемся к функции, описывающей детали ошибки и возвращаем текст на морду
	if resonse.StatusCode != 200 {
		caption := getErrorDescription(resonse.StatusCode)
		ansBody, _ := io.ReadAll(resonse.Body)
		ans := Answer{NormText: caption + "\n" + string(ansBody), RawText: caption + "\n" + string(ansBody)}
		byteAns, _ := json.Marshal(ans)
		return byteAns
	}

	// Читаем из BODY запроса массив байт
	ans, err := io.ReadAll(resonse.Body)

	if err != nil {
		ans := Answer{NormText: err.Error(), RawText: err.Error()}
		byteAns, _ := json.Marshal(ans)
		return byteAns
	}

	// Возвращаем ответ в виде массива байт
	log.Println("Answer from recognition server:", string(ans))
	return ans
}
