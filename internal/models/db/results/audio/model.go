package audio

type DBResult struct {
	FileName string `json:"filename"`
	RawText  string `json:"raw_text"`
	NormText string `json:"norm_text"`
}
