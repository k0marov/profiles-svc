package internal

type Profile struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type ProfileStore interface {
	Get(ID string) (Profile, error)
	Update(ID string, profile Profile) (Profile, error)
}

type ProfileService struct {
	store ProfileStore
}

func NewProfileService(store ProfileStore) *ProfileService {
	return &ProfileService{store: store}
}

func (p ProfileService) Get(ID string) (Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileService) Update(ID string, profile Profile) (Profile, error) {
	//TODO implement me
	panic("implement me")
}
