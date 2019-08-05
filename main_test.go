package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sc-app/app"
	"sc-app/config"
	"sc-app/handler"
	"sc-app/model"
	"testing"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	App   *app.App
	Token string
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	App.Router.ServeHTTP(rr, req)
	return rr
}

func TestMain(m *testing.M) {
	godotenv.Load()
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Cannot connect database %s", err.Error())
	}
	defer db.Close()

	App = &app.App{
		Db:     db,
		Router: chi.NewRouter(),
		Sse:    app.NewSSE(),
	}

	//drop table
	App.Db.DropTableIfExists(&model.User{}, &model.Post{})

	App.Initalize()

	code := m.Run()

	os.Exit(code)
}

func TestRegister(t *testing.T) {
	user := model.User{}
	user.Email = "ashadi.cc@gmail.com"
	user.Password = "123456"
	payload, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusCreated, response.Code)

	json.Unmarshal(response.Body.Bytes(), &user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "", user.Password)

}

func TestLogin(t *testing.T) {
	user := model.User{}
	user.Email = "ashadi.cc@gmail.com"
	user.Password = "123456"
	payload, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code)

	responseToken := handler.ResponseToken{}

	json.Unmarshal(response.Body.Bytes(), &responseToken)

	assert.NotEmpty(t, responseToken.Token)
	assert.NotEmpty(t, responseToken.ExpiredAt)

	Token = responseToken.Token
}
