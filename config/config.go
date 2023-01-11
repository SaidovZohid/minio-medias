package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Minio    Minio
	Postgres PostgresConfig
	HttpPort string
}

type Minio struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func GetConfig(path string) Config {
	godotenv.Load(path + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		HttpPort: conf.GetString("HTTP_PORT"),
		Minio: Minio{
			Endpoint:  conf.GetString("ENDPOINT"),
			AccessKey: conf.GetString("ACCESS_KEY"),
			SecretKey: conf.GetString("SECRET_KEY"),
		},
		Postgres: PostgresConfig{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},
	}
	return cfg
}
