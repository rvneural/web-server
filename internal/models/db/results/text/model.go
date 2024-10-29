package text

type DBResult struct {
	OldText string `json:"old_text"`
	Prompt  string `json:"prompt"`
	NewText string `json:"new_text"`
}
