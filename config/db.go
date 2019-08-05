package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
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
			Host:     os.Getenv("DB_HOST"),
			Charset:  "utf8",
		},
	}

	return &config
}

//ConnectDB connect to database
func ConnectDB() (db *gorm.DB, err error) {
	config := GetConfig()
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Name,
		config.DB.Password)
	db, err = gorm.Open(config.DB.Dialect, dbURI)

	return db, err
}
