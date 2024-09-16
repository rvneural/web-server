package text

type Request struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Text        string `json:"text"`
	Temperature string `json:"temperature"`
}
