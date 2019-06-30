package config

import "os"

//DBConfig struct
type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Port     string
	Charset  string
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
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			Charset:  "utf8",
		},
	}

	return &config
}
