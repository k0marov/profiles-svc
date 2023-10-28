package internal

import (
	"context"
	"encoding/json"
	"net/http"
)
import "github.com/go-chi/chi/v5"

type WebProfileService interface {
	GetOrCreate(ctx context.Context, ID string) (*Profile, error)
	Update(ctx context.Context, profile *Profile) (*Profile, error)
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
		r.Get("/me", s.GetMyProfile)
		r.Patch("/me", s.UpdateProfile)
	})
}

func (s *Server) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	session := GetCaller(r.Context())
	json.NewEncoder(w).Encode(session)
}

func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
