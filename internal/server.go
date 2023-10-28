package internal

import "net/http"
import "github.com/go-chi/chi/v5"

type WebProfileService interface {
	Get(ID string) (Profile, error)
	Update(ID string, profile Profile) (Profile, error)
}

type Server struct {
	svc WebProfileService
	r   chi.Router
}

func NewServer(svc WebProfileService) http.Handler {
	r := chi.NewRouter()
	return &Server{svc, r}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
