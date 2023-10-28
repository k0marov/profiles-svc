package internal

import (
	"context"
	"encoding/json"
	"net/http"
)
import "github.com/go-chi/chi/v5"

type WebProfileService interface {
	GetOrCreate(caller *UserClaims) (*Profile, error)
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
	profile, err := s.svc.GetOrCreate(GetCaller(r.Context()))
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	updated, err := s.svc.GetOrCreate(GetCaller(r.Context()))
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
