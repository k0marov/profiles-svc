package internal

import (
	"encoding/json"
	ory "github.com/ory/client-go"
	"net/http"
)
import "github.com/go-chi/chi/v5"

type WebProfileService interface {
	Get(ID string) (Profile, error)
	Update(ID string, profile Profile) (Profile, error)
}

type Server struct {
	svc WebProfileService
	r   chi.Router
}

func NewServer(svc WebProfileService, client *ory.APIClient) http.Handler {
	srv := &Server{svc, chi.NewRouter()}
	srv.defineEndpoints(client)
	return srv
}

func (s *Server) defineEndpoints(client *ory.APIClient) {
	s.r.Use(NewAuthMiddleware(client))
	s.r.Route("/api/v1/profiles", func(r chi.Router) {
		r.Get("/", s.GetProfile)
		r.Patch("/", s.UpdateProfile)
	})
}

func (s *Server) GetProfile(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r.Context())
	json.NewEncoder(w).Encode(session)
}

func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
