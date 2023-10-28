package internal

type MongoProfileRepository struct {
}

func NewMongoProfileRepository() *MongoProfileRepository {
	return &MongoProfileRepository{}
}

func (m MongoProfileRepository) Get(ID string) (Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoProfileRepository) Create(profile Profile) {
	//TODO implement me
	panic("implement me")
}

func (m MongoProfileRepository) Update(profile Profile) (Profile, error) {
	//TODO implement me
	panic("implement me")
}
