package internal

import (
	"encoding/json"
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

func NewServer(cfg AuthConfig, svc WebProfileService) http.Handler {
	srv := &Server{svc, chi.NewRouter()}
	srv.defineEndpoints(cfg)
	return srv
}

func (s *Server) defineEndpoints(authCfg AuthConfig) {
	s.r.Use(NewAuthMiddleware(authCfg).Middleware())
	s.r.Route("/api/v1/profiles", func(r chi.Router) {
		r.Get("/", s.GetProfile)
		r.Patch("/", s.UpdateProfile)
	})
}

func (s *Server) GetProfile(w http.ResponseWriter, r *http.Request) {
	session := GetUserData(r.Context())
	json.NewEncoder(w).Encode(session)
}

func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
