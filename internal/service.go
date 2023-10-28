package internal

import (
	"context"
	"errors"
	"fmt"
)

type Profile struct {
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
}

var ProfileNotFoundErr = errors.New("profile not found")

type ProfileRepo interface {
	Get(ID string) (*Profile, error)
	Create(profile *Profile) error
	Update(profile *Profile) (*Profile, error)
}

type ProfileService struct {
	repo ProfileRepo
}

func NewProfileService(repo ProfileRepo) *ProfileService {
	return &ProfileService{repo: repo}
}

func (p *ProfileService) GetOrCreate(ctx context.Context, ID string) (*Profile, error) {
	profile, err := p.repo.Get(ID)
	if errors.Is(err, ProfileNotFoundErr) {
		return p.createFirst(GetCaller(ctx))
	}
	if err != nil {
		return nil, fmt.Errorf("while getting profile by id %q from repo: %v", ID, err)
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

func (p *ProfileService) Update(ctx context.Context, profile *Profile) (*Profile, error) {
	profile.ID = GetCaller(ctx).ID
	updated, err := p.repo.Update(profile)
	if err != nil {
		return nil, fmt.Errorf("while updating user profile (id %q) in repo: %v", profile.ID, err)
	}
	return updated, nil
}
