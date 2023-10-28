package internal

import (
	"log"
	"net/http"
)

func InitializeAndStart() {
	repo := NewMongoProfileRepository()
	svc := NewProfileService(repo)
	srv := NewServer(svc)
	log.Fatal(http.ListenAndServe(":8080", srv))
}
