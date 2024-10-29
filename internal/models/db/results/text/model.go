package text

type DBResult struct {
	OldText string `json:"old_text"`
	Prompt  string `json:"prompt"`
	Newtext string `json:"new_text"`
}
