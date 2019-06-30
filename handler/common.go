package handler

import (
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"time"
	"encoding/json"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

var APPLICATION_NAME = "JWT APP"
var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("the secret of kalimdor")

type key int 

//UserInfo key
const UserInfo key = 0


func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

type Message struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func sendMessage(s *sse.Server, path string, m Message) {
	channel := fmt.Sprintf("/events/%s", path)
	payload, _ := json.Marshal(m)
	s.SendMessage(channel, sse.SimpleMessage(string(payload)))
}