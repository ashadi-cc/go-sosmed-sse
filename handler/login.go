package handler 

import (
	"time"
	"github.com/jinzhu/gorm"
	"net/http"
	"sc-app/model"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
)

type responseToken struct {
	Token string `json:"token"`
	ExpiredAt int64 `json:"expired_at"`
}

//Login router
func Login(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	user := model.User{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		RespondError(w, http.StatusBadRequest, "Payload error")
		return
	}
	defer r.Body.Close() 

	if err := db.Where("email = ? AND password = ?", user.Email, user.Password).First(&user).Error; err != nil {
		RespondError(w, http.StatusBadRequest, "user not found")
		return 
	}

	expiredAt := time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix()
	claims := model.MyClaims{
		StandardClaims: jwt.StandardClaims {
			Issuer: APPLICATION_NAME,
			ExpiresAt: expiredAt,
		}, 
		Email: user.Email,
		ID: user.ID,
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
	payload := responseToken {
		Token: signedToken,
		ExpiredAt: expiredAt,
	}
	
	RespondJSON(w, http.StatusOK, payload)

}