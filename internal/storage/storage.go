package storage

import (
	"crypto/rsa"
	"sso/internal/models"
)

type Storage interface {
	Users() Users
	GetRsaKey() (*rsa.PrivateKey, error)
}

type Users interface {
	GetUser(login string) (*models.User, error)
	InsertUser(user *models.User) (string, error)
}

type Migrations interface {
	Migrate(rootPassword string) error
}
