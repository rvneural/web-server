package audio

type Response struct {
	NormText string `json:"normText"`
	RawText  string `json:"rawText"`
	Error    string `json:"error"`
	Details  string `json:"details"`
}
