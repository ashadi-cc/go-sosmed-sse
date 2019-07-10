package main

import (
	"sc-app/config"
	"log"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sc-app/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("no .env file found")
	}
	config := config.GetConfig()

	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Name,
		config.DB.Password)
	db, err := gorm.Open(config.DB.Dialect, dbURI)

	if err != nil {
		log.Fatal("Can't connect to database", err.Error())
	}
	defer db.Close()
	
	app := &app.App{
		Db: db,
		Router: chi.NewRouter(),
		Sse: app.NewSSE(),
	}
	
	app.Run()
}
