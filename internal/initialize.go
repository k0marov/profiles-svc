package internal

import (
	"log"
	"net/http"
)

func InitializeAndStart(cfg AppConfig) {
	repo := NewMongoProfileRepository()
	svc := NewProfileService(repo)
	srv := NewServer(cfg.Auth, svc)
	log.Fatal(http.ListenAndServe(cfg.HTTPServer.Host, srv))
}
