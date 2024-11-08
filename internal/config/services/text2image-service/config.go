package text2image_service

import "os"

var (
	URL = os.Getenv("TEXT2IMAGE_URL")
	KEY = os.Getenv("TEXT2IMAGE_KEY")
)
