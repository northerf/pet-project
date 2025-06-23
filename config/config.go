package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBDsn          string `envconfig:"DB_DSN" default:"postgres://petuser:petpassword@localhost:5432/petprojectdb?sslmode=disable"`
	JwtSecret      []byte `envconfig:"JWT_SECRET" default:"Kc9p1DrGrymYVwGJH5kRatq1g8pK9q2z"`
	Port           string `envconfig:"PORT" default:"8080"`
	LogLevel       string `envconfig:"LOG_LEVEL" default:"debug"`
	MigratitionDir string `envconfig:"MIGRATITIONS_DIR" default:"./migrations"`
}

func Load() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}
