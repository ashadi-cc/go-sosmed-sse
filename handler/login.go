package handler

import (
	"encoding/json"
	"net/http"
	"sc-app/model"
	"sc-app/repo"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type ResponseToken struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}

//Login router
func Login(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	request := model.Login{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&request); err != nil {
		RespondError(w, http.StatusBadRequest, "Payload error")
		return
	}
	defer r.Body.Close()

	user := model.User{}
	user.Email, user.Password = request.Email, request.Password

	repo := &repo.UserRepo{Db: db}

	if err := repo.Login(&user); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	expiredAt := time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix()
	claims := model.MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: expiredAt,
		},
		Email: user.Email,
		ID:    user.ID,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)

	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
	}

	//payload := map[string]string{"token": signedToken}
	payload := ResponseToken{
		Token:     signedToken,
		ExpiredAt: expiredAt,
	}

	RespondJSON(w, http.StatusOK, payload)

}
