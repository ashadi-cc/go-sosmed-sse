package handler

import (
	"encoding/json"
	"net/http"
	"sc-app/model"
	"sc-app/repo"

	"github.com/jinzhu/gorm"
)

//Register New User
func Register(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	user := model.User{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		RespondError(w, http.StatusBadRequest, "Payload error")
		return
	}
	defer r.Body.Close()

	repo := &repo.UserRepo{Db: db}
	if err := repo.Create(&user); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	//reset password
	user.Password = ""
	RespondJSON(w, http.StatusCreated, user)
}
