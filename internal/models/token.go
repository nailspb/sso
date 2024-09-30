package models

import "time"

type Token struct {
	Id         string
	User       string
	CreateAt   time.Time
	ValidUntil time.Time
}
