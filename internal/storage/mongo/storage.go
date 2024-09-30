package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"sso/internal/models"
	"sso/internal/services"
	"sso/internal/storage"
	"sso/pkg/helpers/errorHelper"
	"sso/pkg/helpers/slogHelper"
	"time"
)

const (
	ErrorCreateClient            = "Error on create mongoDB client"
	ErrorCreateMigrationInstance = "Error on create migration instance"
	ErrorCloseConnections        = "Error on close client connections"
	ErrorCreateRoot              = "Error on insert root user into DB"
	ErrorInternal                = "Internal mongo server error"
)

type Storage struct {
	client *mongo.Client
	db     *mongo.Database
}

type Config struct {
	Server       string
	User         string
	Password     string
	Database     string
	RootPassword string
}

func New(logger *slog.Logger, config Config) (*Storage, error) {
	const operation = "internal.storage.mongo.new()"
	log := slogHelper.AddOperation(logger, operation)
	logger.Info("Connecting to MongoDB and init it...")
	//
	//mongo client
	//
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", config.User, config.Password, config.Server)).
		SetMaxConnecting(50). //максимально количество одновременных соединений
		SetMaxPoolSize(100),  //размер пула соединений
	)
	if err != nil {
		return nil, errorHelper.WrapError(operation, ErrorCreateClient, err)
	}
	//
	//migrations
	//
	driver, _ := mongodb.WithInstance(client, &mongodb.Config{
		DatabaseName: "sso",
	})
	m, err := migrate.NewWithDatabaseInstance("file://migrations/mongo", "mongo", driver)
	if err != nil {
		return nil, errorHelper.WrapError(operation, ErrorCreateMigrationInstance, err)
	}
	err = m.Up()
	if err != nil {
		log.Warn("Migration", slogHelper.GetErrAttr(err))
	}
	s := &Storage{
		client: client,
		db:     client.Database(config.Database),
	}
	//
	//add root user
	//
	_, err = s.Users().GetUser("root")
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			root := &models.User{
				Login:    "root",
				Password: config.RootPassword,
			}
			insertId, err := services.Users(s).Add(root)
			if err != nil {
				log.Error(ErrorCreateRoot, slogHelper.GetErrAttr(err))
			} else {
				logger.Info("------------------------------------------------------------------------------------------")
				logger.Info("Successfully created user",
					slog.String("uid", insertId),
					slog.String("login", "root"),
					slog.String("password", config.RootPassword),
				)
				logger.Info("------------------------------------------------------------------------------------------")
			}
		} else {
			logger.Error(ErrorInternal, slogHelper.GetErrAttr(err))
		}
	}
	return s, nil
}

func (s *Storage) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := s.client.Ping(ctx, nil)
	if err != nil {
		return false
	}
	return true
}

func (s *Storage) Shutdown(ctx context.Context) error {
	const operation = "internal.storage.mongo.new()"
	err := s.client.Disconnect(ctx)
	if err != nil {
		return errorHelper.WrapError(operation, ErrorCloseConnections, err)
	}
	return nil
}

func (s *Storage) Users() storage.Users {
	return &Users{
		db: s.db,
	}
}
