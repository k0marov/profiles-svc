package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrProfileNotFound = ClientError{
	DisplayMessage: "profile not found",
	HTTPCode:       http.StatusNotFound,
}

type ProfileRepo interface {
	Get(ID string) (*Profile, error)
	Create(profile *Profile) error
	Update(ID string, profile *ProfileUpdatable) (*Profile, error)
}

type ProfileService struct {
	repo ProfileRepo
}

func NewProfileService(repo ProfileRepo) *ProfileService {
	return &ProfileService{repo: repo}
}

func (p *ProfileService) GetOrCreate(caller *UserClaims) (*Profile, error) {
	profile, err := p.repo.Get(caller.ID)
	if errors.Is(err, ErrProfileNotFound) {
		return p.createFirst(caller)
	}
	if err != nil {
		return nil, fmt.Errorf("while getting profile by id %q from repo: %v", caller.ID, err)
	}
	return profile, nil
}

func (p *ProfileService) createFirst(user *UserClaims) (*Profile, error) {
	email, ok := user.Traits["email"].(string)
	if !ok {
		email = ""
	}
	profile := &Profile{
		ID:    user.ID,
		Email: email,
	}
	err := p.repo.Create(profile)
	if err != nil {
		return nil, fmt.Errorf("while creating user profile (id %q) on first request: %v", user.ID, err)
	}
	return profile, nil
}

func (p *ProfileService) Update(ctx context.Context, upd *ProfileUpdatable) (*Profile, error) {
	profile, err := p.repo.Update(GetCaller(ctx).ID, upd)
	if err != nil {
		return nil, fmt.Errorf("while updating user profile in repo: %v", err)
	}
	return profile, nil
}
