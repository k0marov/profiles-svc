package internal

import (
	"log"
	"net/http"
)

func InitializeAndStart(cfg AppConfig) {
	repo, closeRepo := NewMongoProfileRepository(cfg.Mongo)
	defer closeRepo()
	svc := NewProfileService(repo)
	srv := NewServer(cfg.Auth, svc)
	log.Print(http.ListenAndServe(cfg.HTTPServer.Host, srv))
}
