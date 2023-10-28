package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoProfileRepository struct {
	col *mongo.Collection
}

func NewMongoProfileRepository(cfg MongoConfig) (repo *MongoProfileRepository, close func()) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		log.Fatalf("unable to connect to mongo db: %v", err)
	}
	log.Printf("connected to mongodb at %q", cfg.URI)
	close = func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("unable to close mongo db connection: %v", err)
		}
	}
	profilesCol := client.Database("profiles-svc").Collection("profiles")
	return &MongoProfileRepository{
		profilesCol,
	}, close
}

func (m *MongoProfileRepository) Get(ID string) (*Profile, error) {
	result := m.col.FindOne(context.TODO(), bson.D{{"_id", ID}}) // Profile{ID: ID})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("while finding profile with ID=%q: %v", ID, err)
	}
	var profile Profile
	if err := result.Decode(&profile); err != nil {
		return nil, fmt.Errorf("unable to decode profile from db: %v", err)
	}
	return &profile, nil
}

func (m *MongoProfileRepository) Create(profile *Profile) error {
	_, err := m.col.InsertOne(context.TODO(), profile)
	if err != nil {
		return fmt.Errorf("while inserting profile into db: %v", err)
	}
	return nil
}

func (m *MongoProfileRepository) Update(ID string, upd *ProfileUpdatable) (*Profile, error) {
	profile, err := m.Get(ID)
	if err != nil {
		return nil, fmt.Errorf("while getting profile profile: %v", err)
	}
	profile.ProfileUpdatable = updateProfile(profile.ProfileUpdatable, *upd)

	result := m.col.FindOneAndReplace(context.TODO(), bson.D{{"_id", ID}}, profile)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("while updating profile in db: %v", err)
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
