package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

var Config ProjectConfig

type ProjectConfig struct {
	Storage  string   `yaml:"storage" env:"STORAGE" env-default:"inmemory"`
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
}

type Server struct {
	Host string `yaml:"host" env:"SRV_HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"SRV_PORT" env-default:"8000"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"PG_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"PG_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-default:"postgres"`
	Database string `yaml:"database" env:"PG_DB" env-default:"posts"`
	SslMode  string `yaml:"sslmode" env:"PG_SSL_MODE" env-default:"disable"`
}

func Read() error {
	err := cleanenv.ReadConfig("config/config.yaml", &Config)
	if err != nil {
		return fmt.Errorf("error while reading application configuration: %w", err)
	}

	err = cleanenv.ReadEnv(&Config)
	if err != nil {
		return fmt.Errorf("error while reading environment variables: %w", err)
	}

	return nil
}
