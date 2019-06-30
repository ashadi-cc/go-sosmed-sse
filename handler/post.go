package handler

import (
	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-chi/chi"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"sc-app/model"
	jwt "github.com/dgrijalva/jwt-go"
)

//CreatePost new post
func CreatePost(w http.ResponseWriter, r *http.Request, db *gorm.DB, s *sse.Server) {
	userInfo := r.Context().Value(UserInfo).(jwt.MapClaims)
	post := model.Post{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&post); err != nil {
		RespondError(w, http.StatusBadRequest, "Payload error")
		return
	}
	defer r.Body.Close() 

	userID := fmt.Sprintf("%v", userInfo["id"])
	id, _:= strconv.ParseUint(userID, 10, 32)
	post.UserId = uint(id)
	post.ID = 0
	if err := db.Save(&post).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, "internal server error")
		return 
	}

	RespondJSON(w, http.StatusCreated, post)

	sendMessage(s, "post", Message{
		Status : "created",
		Data: post,
	})
}

//GetAllPost get all post
func GetAllPost(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	posts := []model.Post{}

	if err := db.Order("id DESC").Find(&posts).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, "internal server error")
		return 		
	}

	RespondJSON(w, http.StatusOK, posts)
}

//GetPost get single post
func GetPost(w http.ResponseWriter, r *http.Request, db *gorm.DB, s *sse.Server) {
	postID := chi.URLParam(r, "PostID")
	post := model.Post{}

	if err := db.First(&post, postID).Error; err != nil {
		RespondError(w, http.StatusNotFound, "post not found")
		return
	}

	RespondJSON(w, http.StatusOK, post)
}

//UpdatePost post
func UpdatePost(w http.ResponseWriter, r *http.Request, db *gorm.DB, s *sse.Server) {
	postID := chi.URLParam(r, "PostID")
	post := model.Post{}

	if err := db.First(&post, postID).Error; err != nil {
		RespondError(w, http.StatusNotFound, "post not found")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&post); err != nil {
		RespondError(w, http.StatusBadRequest, "Payload error")
		return
	}
	defer r.Body.Close() 

	//reset postID
	id, _:= strconv.ParseUint(postID, 10, 32)
	post.ID = uint(id)
	if err := db.Save(&post).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, "internal server error")
		return 
	}

	RespondJSON(w, http.StatusOK, post)

	sendMessage(s, "post", Message{
		Status : "updated",
		Data: post,
	})
}

//DeletePost delete single post
func DeletePost(w http.ResponseWriter, r *http.Request, db *gorm.DB, s *sse.Server) {
	postID := chi.URLParam(r, "PostID")
	post := model.Post{}

	if err := db.First(&post, postID).Error; err != nil {
		RespondError(w, http.StatusNotFound, "post not found")
		return
	}

	if err := db.Exec("DELETE FROM posts WHERE id = ?", postID).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	RespondJSON(w, http.StatusOK, post)

	sendMessage(s, "post", Message{
		Status : "deleted",
		Data: post,
	})
}