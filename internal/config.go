package internal

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type AuthConfig struct {
	OryProxyPort string `default:"4000"`
	OryProxyURL  string `default:"http://localhost:4000/.ory"`
}

type HTTPServerConfig struct {
	Host string `default:"localhost:3000"`
}

type AppConfig struct {
	Auth       AuthConfig
	HTTPServer HTTPServerConfig
}

func ReadConfigFromEnv() AppConfig {
	var cfg AppConfig
	err := envconfig.Process("profiles", &cfg)
	if err != nil {
		log.Panicf("while parsing app config from env: %v", err)
	}
	return cfg
}
