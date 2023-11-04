package internal

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrProfileNotFound = &ClientError{
	DisplayMessage: "profile not found",
	HTTPCode:       http.StatusNotFound,
}

type ProfileRepo interface {
	Get(ID string) (*Profile, error)
	Create(profile *Profile) error
	Replace(ID string, profile *Profile) error
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
		return nil, fmt.Errorf("while getting profile by id %q from repo: %w", caller.ID, err)
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
		return nil, fmt.Errorf("while creating user profile (id %q) on first request: %w", user.ID, err)
	}
	return profile, nil
}

func (p *ProfileService) Get(userID string) (*Profile, error) {
	profile, err := p.repo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("while getting profile %q from repo: %v", userID, err)
	}
	return profile, nil
}

func (p *ProfileService) Update(caller *UserClaims, upd *ProfileUpdatable) (*Profile, error) {
	profile, err := p.repo.Get(caller.ID)
	if err != nil {
		return nil, fmt.Errorf("while getting profile: %w", err)
	}
	profile.ProfileUpdatable = updateProfile(profile.ProfileUpdatable, *upd)

	err = p.repo.Replace(caller.ID, profile)
	if err != nil {
		return nil, fmt.Errorf("while replacing user profile in repo: %w", err)
	}
	return profile, nil
}

func updateProfile(current, upd ProfileUpdatable) ProfileUpdatable {
	if upd.Name != nil {
		current.Name = upd.Name
	}
	if upd.Age != nil {
		current.Age = upd.Age
	}
	return current
}
