package mongo

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sso/internal/models"
	"sso/pkg/helpers/errorHelper"
	"time"
)

type Keys struct {
	db *mongo.Database
}

const (
	ErrorGenRsaKey       = "error on generating rsa key"
	ErrorKeyDecode       = "Error on decode key document"
	ErrorSaveKey         = "error on save private key to DB"
	ErrorParsePrivateKey = "error on parsing private key"
	ErrorCastPrivateKey  = "error on casting private key"
)

func (s *Storage) generateKey() (*rsa.PrivateKey, error) {
	const operation = "internal.storage.mongo.generateKey()"
	//generate rsa key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errorHelper.WrapError(operation, ErrorGenRsaKey, err)
	}
	//encode keys to PEM format
	/*privatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)*/
	//write key to DataBase
	_, err = s.db.Collection("Keys").InsertOne(context.TODO(), bson.D{
		{"PEM", x509.MarshalPKCS1PrivateKey(privateKey)},
		{"Exp", time.Now().Add(time.Hour * 2)},
	})
	if err != nil {
		return nil, errorHelper.WrapError(operation, ErrorSaveKey, err)
	}
	//return key
	return privateKey, nil
}

func (s *Storage) GetRsaKey() (*rsa.PrivateKey, error) {
	const operation = "internal.storage.mongo.GetKey()"
	//find last key
	m := models.Key{}
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	find := s.db.Collection("Keys").FindOne(context.TODO(), bson.M{}, opts)
	//if key not found
	if err := find.Err(); err != nil {
		return s.generateKey()
	}
	//
	// if key found in DB
	//

	//decode data
	if err := find.Decode(&m); err != nil {
		return nil, errorHelper.WrapError(operation, ErrorKeyDecode, err)
	}
	//check exp
	if time.Now().After(m.Exp) {
		return s.generateKey()
	}
	//decode PEM format
	//privateKey, err := ssh.ParseRawPrivateKey([]byte(m.PEM))
	privateKey, err := x509.ParsePKCS1PrivateKey(m.PEM)
	if err != nil {
		return nil, errorHelper.WrapError(operation, ErrorParsePrivateKey, err)
	}
	//cast to rsa key
	/*ret, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errorHelper.WrapError(operation, ErrorCastPrivateKey, err)
	}*/
	return privateKey, nil
}
