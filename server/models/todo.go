package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title"`
	Body  string             `json:"body" bson:"body"`
	Done  bool               `json:"done" bson:"done"`
}
