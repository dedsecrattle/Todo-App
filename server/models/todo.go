package models

type Todo struct {
	ID    int    `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Body  string `json:"body" bson:"body"`
	Done  bool   `json:"done" bson:"done"`
}
