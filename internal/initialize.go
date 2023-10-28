package internal

import (
	"log"
	"net/http"
)

func InitializeAndStart() {
	srv := NewServer()
	log.Fatal(http.ListenAndServe(":8080", srv))
}
