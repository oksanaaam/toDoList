package model

type ToDo struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
