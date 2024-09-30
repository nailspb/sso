package models

import "time"

type Key struct {
	Id  string `bson:"_id, omitempty"`
	PEM []byte
	Exp time.Time
}
