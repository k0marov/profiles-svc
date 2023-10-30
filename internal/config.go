package internal

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type HTTPServerConfig struct {
	Host string `default:"127.0.0.1:8001"`
}

type MongoConfig struct {
	URI string `default:"mongodb://127.0.0.1:27017"`
}

type AppConfig struct {
	HTTPServer HTTPServerConfig
	Mongo      MongoConfig
}

func ReadConfigFromEnv() AppConfig {
	var cfg AppConfig
	err := envconfig.Process("profiles", &cfg)
	if err != nil {
		log.Panicf("while parsing app config from env: %v", err)
	}
	return cfg
}
