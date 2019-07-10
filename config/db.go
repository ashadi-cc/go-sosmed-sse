package config

import (
	"os"
)

//DBConfig struct
type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Port     string
	Charset  string
	Host     string
}

//Config struct
type Config struct {
	DB *DBConfig
}

//GetConfig function
func GetConfig() *Config {
	config := Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_DB"),
			Port:     os.Getenv("DB_PORT"),
			Host:	  os.Getenv("DB_HOST"),
			Charset:  "utf8",
		},
	}

	return &config
}
