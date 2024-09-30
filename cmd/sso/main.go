package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sso/internal/config/env"
	"sso/internal/http/handlers"
	"sso/internal/storage/mongo"
	"sso/pkg/helpers/slogHelper"
	"sso/pkg/http/middleware"
	"sso/pkg/http/routing"
	"time"
)

func main() {
	config := env.New()
	log := slogHelper.NewLogger(config.DebugLevel)
	log.Info("start SSO service", slog.String("env", config.DebugLevel))
	log.Debug("debug messages are enabled")

	storage, err := mongo.New(log, mongo.Config{
		Server:       config.Db.Server,
		User:         config.Db.User,
		Password:     config.Db.Password,
		Database:     config.Db.Database,
		RootPassword: config.RootPassword,
	})
	if err != nil {
		log.Error("failed init storage", slogHelper.GetErrAttr(err))
		os.Exit(1)
	}
	//storage.InitDatabase(log, config)

	/*storage, err := sqlite.New(cfg.Sqlite.Path, cfg.Rsa.Folder)
	if err != nil {
		log.Error("failed init storage", slogHelper.GetErrAttr(err))
		os.Exit(1)
	}
	//do data migrations
	if err := storage.Migrations().Migrate(cfg.Init.Password); err != nil {
		log.Error("failed migrate database", slogHelper.GetErrAttr(err))
	}*/

	//configure routes
	routes := routing.New().
		Handle("POST /{$}", handlers.Auth(log, storage)).
		Handle("GET /status", handlers.Status(log)).
		Handle("GET /key", handlers.Key(log, storage)).
		Handle("POST /check", handlers.Check(log, storage)).
		UseMiddleware(middleware.Logging(log)).
		UseMiddleware(middleware.Recovery(log, config.DebugLevel)).
		UseMiddleware(middleware.RequestId(log))

	//configure http server
	srv := &http.Server{
		Addr:           config.Server.Address,
		Handler:        routes,
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: 1 * 1024 * 1024, //1Mb
	}

	//start http server
	log.Info("starting http server", slog.String("address", config.Server.Address))
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error("Failed start server", slogHelper.GetErrAttr(err))
		}
	}()

	//graceful shutdown
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)
	sig := <-channel
	log.Info("Stop signal received", slog.String("signal", sig.String()))
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	if err := storage.Shutdown(ctx); err != nil {
		log.Error("Storage Graceful shutdown failed", slogHelper.GetErrAttr(err))
	}
	log.Info("Server shutdown successfully")

}
