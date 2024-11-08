package app

import "os"

const (
	HTTP_PORT  = ":80"
	HTTPS_PORT = ":443"
	DOMAIN     = "neuron-nexus.ru"
)

var (
	SESSION_SECRET = os.Getenv("SESSION_SECRET")
	ROLLBAR_TOKEN  = os.Getenv("ROLLBAR_TOKEN")
	ADMIN_LOGIN    = os.Getenv("ADMIN_LOGIN")
	ADMIN_PASSWORD = os.Getenv("ADMIN_PASSWORD")
)
