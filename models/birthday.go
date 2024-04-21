package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BirthDay struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"` // Maps the MongoDB _id to JSON "id"
	Name  string             `json:"name" bson:"name"`
	Day   int                `json:"day" bson:"day"`
	Month int                `json:"month" bson:"month"`
	Year  int                `json:"year" bson:"year"`
}
