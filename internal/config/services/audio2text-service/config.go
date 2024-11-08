package audio2text_service

import "os"

var (
	URL = os.Getenv("AUDIO2TEXT_URL")
	KEY = os.Getenv("AUDIO2TEXT_KEY")
)
