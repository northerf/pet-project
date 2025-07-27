package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBDsn          string `envconfig:"DB_DSN" default:""`
	JwtSecret      []byte `envconfig:"JWT_SECRET" default:""`
	Port           string `envconfig:"PORT" default:""`
	LogLevel       string `envconfig:""`
	MigratitionDir string `envconfig:""`
}

func Load() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}
