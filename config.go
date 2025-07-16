package main

import (
	"os"

	"github.com/yureien/animeserver/database"
	"github.com/yureien/animeserver/filehandler"
	"github.com/yureien/animeserver/server"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database    database.DatabaseConfig `yaml:"database"`
	Server      server.Config           `yaml:"server"`
	FileHandler filehandler.Config      `yaml:"file_handler"`
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
