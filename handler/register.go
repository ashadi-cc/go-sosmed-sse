package handler 

import (
	"encoding/json"
	"sc-app/model"
	"github.com/jinzhu/gorm"
	"net/http"
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

	if err := db.Where("email = ?", user.Email).First(&user).Error; err == nil && user.ID != 0 {
		RespondError(w, http.StatusBadRequest, "user exsist")
		return 
	}

	if err := db.Save(&user).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, "internal server error")
		return 
	}

	RespondJSON(w, http.StatusCreated, user)
}