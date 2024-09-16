package text

type Response struct {
	NewText string `json:"newText"`
	OldText string `json:"oldText"`
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
