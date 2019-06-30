package mdware

import (
	"log"
	"context"
	"fmt"
	"sc-app/handler"
	"strings"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

//NoAuth router with no need auth
var NoAuth = map[string]string {"/register": "noauth", "/login": "noauth"}

func AddNoAuthRoute(path string) {
	NoAuth[path] = "noauth"
}

//JWTAuth auth JWT
func JWTAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if _, ok := NoAuth[path]; ok {
			log.Println("No auth", path)
			next.ServeHTTP(w, r)
			return
		} 

		//check token here
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		if tokenString == "" {
			tokenString = r.URL.Query().Get("token")
		}

		if tokenString == "" {
			handler.RespondError(w, http.StatusBadRequest, "Invalid token")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != handler.JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("Signing method invalid")
			}
		
			return handler.JWT_SIGNATURE_KEY, nil
		})

		if  err != nil {
			handler.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			handler.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), handler.UserInfo, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
