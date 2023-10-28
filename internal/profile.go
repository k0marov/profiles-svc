package internal

type Profile struct {
	ID    string `json:"id" bson:"_id"`
	Email string `json:"email" bson:"email"`
	ProfileUpdatable
}

type ProfileUpdatable struct {
	Name *string `json:"name" bson:"name"`
	Age  *int    `json:"age" bson:"age"`
}
