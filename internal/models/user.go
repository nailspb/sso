package models

type User struct {
	Id       string `bson:"_id, omitempty"`
	Login    string
	Password string
}
