package main

import (
	"log"
	"sc-app/app"
	"sc-app/config"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, load env variables from local!")
	}
	db, err := config.ConnectDB()

	if err != nil {
		log.Fatal("Can't connect to database", err.Error())
	}
	defer db.Close()

	app := &app.App{
		Db:     db,
		Router: chi.NewRouter(),
		Sse:    app.NewSSE(),
	}

	app.Run()
}
