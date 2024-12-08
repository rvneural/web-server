package text

type Request struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Text        string `json:"text"`
	Temperature string `json:"temperature"`
	OperationId string `json:"operation_id"`
	UserID      int    `json:"user_id" xml:"user_id" form:"user_id"`
}
