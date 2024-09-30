package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sso/internal/models"
	"sso/pkg/helpers/errorHelper"
)

type Users struct {
	db *mongo.Database
}

const (
	ErrorUserNotFound = "User not found"
	ErrorUserDecode   = "Error on decode user document"
	ErrorInsertUser   = "Error on insert user document"
)

func (s *Users) GetUser(login string) (*models.User, error) {
	const operation = "internal.storage.mongo.GetUser()"
	find := s.db.Collection("Users").FindOne(context.TODO(), bson.M{"login": login})
	if err := find.Err(); err != nil {
		return nil, errorHelper.WrapError(operation, ErrorUserNotFound, err)
	}
	user := models.User{}
	err := find.Decode(&user)
	if err != nil {
		return nil, errorHelper.WrapError(operation, ErrorUserDecode, err)
	}
	return &user, nil
}

func (s *Users) InsertUser(user *models.User) (string, error) {
	const operation = "internal.storage.mongo.InsertUser()"
	res, err := s.db.Collection("Users").InsertOne(context.TODO(), bson.D{
		{"login", user.Login},
		{"password", user.Password},
	})
	if err != nil {
		return "", errorHelper.WrapError(operation, ErrorInsertUser, err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
