package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Country struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Age        int                `bson:"age" json:"age"`
	Population int                `bson:"population" json:"population"`
	Continent  string             `bson:"continent" json:"continent"`
}
