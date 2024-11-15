package app

import "os"

type Config struct {
	SESSION_SECRET string
	ROLLBAR_TOKEN  string
	ADMIN_LOGIN    string
	ADMIN_PASSWORD string
	HTTP_PORT      string
	HTTPS_PORT     string
	DOMAIN         string
}

func Init() *Config {
	return &Config{
		HTTP_PORT:      ":80",
		HTTPS_PORT:     ":443",
		DOMAIN:         "neuron-nexus.ru",
		SESSION_SECRET: os.Getenv("SESSION_SECRET"),
		ROLLBAR_TOKEN:  os.Getenv("ROLLBAR_TOKEN"),
		ADMIN_LOGIN:    os.Getenv("ADMIN_LOGIN"),
		ADMIN_PASSWORD: os.Getenv("ADMIN_PASSWORD"),
	}
}
