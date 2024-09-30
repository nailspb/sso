package config

import "time"

type Config struct {
	DebugLevel   string
	RootPassword string
	Server       ServerConfig
	Db           DbConfig
}

type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DbConfig struct {
	Server   string
	User     string
	Password string
	Database string
}
