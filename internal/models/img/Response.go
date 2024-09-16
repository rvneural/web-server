package img

type Response struct {
	Error   string `json:"error,omitempty"`
	Details string `json:"details"`
	Image   struct {
		B64String string `json:"b64String"`
		Seed      string `json:"seed"`
	} `json:"image"`
}
