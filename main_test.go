package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	App.Db.DropTableIfExists(&model.Post{}, &model.User{})

	App.Initalize()

	code := m.Run()

	os.Exit(code)
}

func TestRegister(t *testing.T) {
	user := model.Login{}
	user.Email = "ashadi.cc@gmail.com"
	user.Password = "123456"
	payload, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusCreated, response.Code, string(response.Body.Bytes()))

	userCreated := model.User{}
	json.Unmarshal(response.Body.Bytes(), &userCreated)
	assert.Equal(t, uint(1), userCreated.ID)
	assert.Equal(t, "", userCreated.Password)

}

func TestLogin(t *testing.T) {
	user := model.Login{}
	user.Email = "ashadi.cc@gmail.com"
	user.Password = "123456"
	payload, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code, string(response.Body.Bytes()))

	responseToken := handler.ResponseToken{}

	json.Unmarshal(response.Body.Bytes(), &responseToken)

	assert.NotEmpty(t, responseToken.Token)
	assert.NotEmpty(t, responseToken.ExpiredAt)

	Token = responseToken.Token
}

func TestCreatePost(t *testing.T) {
	post := model.Post{}
	post.Title, post.Body = "Sumanto", "Sujono"

	payload, _ := json.Marshal(post)

	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(payload))

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusCreated, response.Code, string(response.Body.Bytes()))
}

func TestGetAllPost(t *testing.T) {
	req, err := http.NewRequest("GET", "/post", nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code, string(response.Body.Bytes()))

	assert.NotEqual(t, "[]", string(response.Body.Bytes()))
}

func TestFindPost(t *testing.T) {
	req, err := http.NewRequest("GET", "/post/1", nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code, string(response.Body.Bytes()))

	assert.NotEqual(t, "[]", string(response.Body.Bytes()))
}

func TestUpdatePost(t *testing.T) {
	post := model.Post{}
	post.Title, post.Body = "mata", "hati"

	payload, _ := json.Marshal(post)

	req, err := http.NewRequest("PUT", "/post/1", bytes.NewBuffer(payload))

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code, string(response.Body.Bytes()))
}

func TestDeletePost(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/post/1", nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))

	assert.Equal(t, nil, err)

	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code, string(response.Body.Bytes()))
}
