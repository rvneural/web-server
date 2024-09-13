package scripts

import "strconv"

// Возращает описание ошибки, которая может прийти с Main Server
func getErrorDescription(statusCode int) string {
	switch statusCode {
	case 404:
		return "Error 404 — Не удалось обратиться к серверу расшифровки"
	case 204:
		return "Error 204 — Вы не передали содержимое файла"
	case 400:
		return "Error 400 — Ошибка в сигнатуре запроса"
	case 415:
		return "Error 415 — Тип файла не поддерживается"
	case 521:
		return "Error 521 — Сервер расшифровки недоступн"
	}
	return "HTTP Error " + strconv.Itoa(statusCode)
}
