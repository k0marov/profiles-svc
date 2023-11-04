package internal

import (
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "gitlab.com/samkomarov/profiles-svc.git/docs"
	"net/http"
)
import "github.com/go-chi/chi/v5"

type WebProfileService interface {
	GetOrCreate(caller *UserClaims) (*Profile, error)
	Get(userID string) (*Profile, error)
	Update(caller *UserClaims, upd *ProfileUpdatable) (*Profile, error)
}

type Server struct {
	svc WebProfileService
	r   chi.Router
}

func NewServer(svc WebProfileService) http.Handler {
	srv := &Server{svc, chi.NewRouter()}
	srv.defineEndpoints()
	return srv
}

//	@title			profiles-svc
//	@version		1.0
//	@description	A microservice for handling user profiles

//	@contact.name	Sam Komarov
//	@contact.url	github.com/k0marov
//	@contact.email	sam@skomarov.com

// @host		localhost:8080
// @schemes     https http
func (s *Server) defineEndpoints() {
	s.r.Get("/swagger/*", httpSwagger.WrapHandler)
	s.r.Route("/api/v1/profiles", func(r chi.Router) {
		r.Use(AuthMiddleware())
		r.Get("/me", s.GetMyProfile)
		r.Get("/{id}", s.GetOtherProfile)
		r.Patch("/me", s.UpdateProfile)
	})
}

// GetMyProfile godoc
//
//		@Summary		Get caller's profile
//		@Description	Get profile of the caller if it has been created.
//	    @Description 	If profile was not yet created, create it.
//		@Tags			profiles
//		@Produce		json
//		@Success		200	{object}	[]Profile
//		@Router			/api/v1/profiles/me [get]
func (s *Server) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := s.svc.GetOrCreate(GetCaller(r.Context()))
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

// GetOtherProfile godoc
//
//		@Summary		Get profile by user id
//		@Description	Get profile by user id. Returns 404 if profile does not exist.
//		@Tags			profiles
//		@Produce		json
//		@Success		200	{object}	Profile
//	    @Failure        404 {object}    ClientError
//		@Param			id	path		string	true	"ID of the user for which you want to get its profile."
//		@Router			/api/v1/profiles/{id} [get]
func (s *Server) GetOtherProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := s.svc.Get(chi.URLParam(r, "id"))
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile godoc
//
//	@Summary		Update profile of the caller
//	@Description	Update profile of the caller, only updating the specified fields.
//	@Tags			profiles
//	@Accept 		json
//	@Param			account	body		ProfileUpdatable	true	"fields to update"
//	@Produce		json
//	@Success		200	{object}	Profile
//	@Router			/api/v1/profiles/me [patch]
func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var upd ProfileUpdatable
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	updated, err := s.svc.Update(GetCaller(r.Context()), &upd)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
