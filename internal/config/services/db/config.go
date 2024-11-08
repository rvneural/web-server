package db

import "os"

var (
	URL = os.Getenv("DB_URL")
)
