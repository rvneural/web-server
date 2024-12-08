package img

type Request struct {
	Prompt      string `json:"prompt"`
	Seed        string `json:"seed"`
	WidthRatio  string `json:"widthRatio"`
	HeightRatio string `json:"heightRatio"`
	OperationId string `json:"operation_id"`
	UserID      int    `json:"user_id" xml:"user_id" form:"user_id"`
}
