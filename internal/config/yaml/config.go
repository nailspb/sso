package yaml

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
)

//TODO:  need refactoring

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-default:"prod"`
	Server struct {
		Address string `yaml:"address" env-default:"localhost:8085"`
	} `yaml:"server" jsonHelper:"server"`
	Init struct {
		Password string `yaml:"password" env-required:"true"`
	} `yaml:"init"`
	Sqlite struct {
		Path string `yaml:"path" env-default:"./db.db"`
	} `yaml:"sqlite"`
	Rsa struct {
		Folder string `yaml:"folder" env-default:"./keys"`
	} `yaml:"rsa"`
}

func Load() *Config {
	confPath := flag.String("config", configDefaultPath(), "config file path")
	flag.Parse()
	if _, err := os.Stat(*confPath); err != nil {
		confPath2 := appPath() + "/" + *confPath
		if _, err := os.Stat(confPath2); err != nil {
			log.Fatal("bad config file path: ", *confPath)
		}
		confPath = &confPath2
	}

	config := &Config{}
	if err := cleanenv.ReadConfig(*confPath, config); err != nil {
		log.Fatal("cannot parse config file ", *confPath, ": ", err)
	}
	return config
}

func configDefaultPath() string {
	return appPath() + "/configs/local.yaml"
}

func appPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)
	return path
}
