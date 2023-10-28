package internal

import "net/http"

type Server struct {
}

func NewServer() http.Handler {
	return &Server{}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
