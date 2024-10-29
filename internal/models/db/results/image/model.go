package image

type DBResult struct {
	Prompt    string `json:"prompt"`
	Seed      string `json:"seed"`
	B64string string `json:"image"`
}
