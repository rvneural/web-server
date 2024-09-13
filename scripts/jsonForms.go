package scripts

// Стркутура запроса на Main Server
type Request struct {
	FileData  []byte   `json:"fileData"`
	Languages []string `json:"languages"`
	FileType  string   `json:"fileType"`
	Dialog    bool     `json:"dialog"`
}

// Структура ответа от Main Server, которая в дальнейшем парсится в JS
// и отправляется на мордду
type Answer struct {
	NormText string `json:"normText"`
	RawText  string `json:"rawText"`
}
