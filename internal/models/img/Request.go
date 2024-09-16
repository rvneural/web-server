package img

type Request struct {
	Prompt      string `json:"prompt"`
	Seed        string `json:"seed"`
	WidthRatio  string `json:"widthRatio"`
	HeightRatio string `json:"heightRatio"`
}
