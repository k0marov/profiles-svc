package internal

import (
	ory "github.com/ory/client-go"
	"log"
	"net/http"
)

func InitializeAndStart(cfg AppConfig) {
	oryCfg := ory.NewConfiguration()
	oryCfg.Servers = ory.ServerConfigurations{{URL: cfg.Auth.OryProxyURL}}

	repo := NewMongoProfileRepository()
	svc := NewProfileService(repo)
	srv := NewServer(svc, ory.NewAPIClient(oryCfg))
	log.Fatal(http.ListenAndServe(cfg.HTTPServer.Host, srv))
}
