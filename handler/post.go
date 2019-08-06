package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sc-app/model"
	"sc-app/repo"
	"strconv"

	"github.com/alexandrevicenzi/go-sse"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
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
	id, _ := strconv.ParseUint(userID, 10, 32)
	post.UserID = uint(id)
	post.ID = 0
	repo := &repo.PostRepo{Db: db}
	if err := repo.Create(&post); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, post)

	sendMessage(s, "post", Message{
		Status: "created",
		Data:   post,
	})
}

//GetAllPost get all post
func GetAllPost(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	repo := &repo.PostRepo{Db: db}
	posts, err := repo.All()

	if err != nil {
		RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	RespondJSON(w, http.StatusOK, posts)
}

//GetPost get single post
func GetPost(w http.ResponseWriter, r *http.Request, db *gorm.DB, s *sse.Server) {
	postID := chi.URLParam(r, "PostID")
	repo := &repo.PostRepo{Db: db}

	id, _ := strconv.ParseUint(postID, 10, 32)
	post, err := repo.FindByID(uint(id))

	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, post)
}

//UpdatePost post
func UpdatePost(w http.ResponseWriter, r *http.Request, db *gorm.DB, s *sse.Server) {
	postID := chi.URLParam(r, "PostID")
	post := model.Post{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&post); err != nil {
		RespondError(w, http.StatusBadRequest, "Payload error")
		return
	}
	defer r.Body.Close()

	repo := &repo.PostRepo{Db: db}
	id, _ := strconv.ParseUint(postID, 10, 32)
	post.ID = uint(id)

	userInfo := r.Context().Value(UserInfo).(jwt.MapClaims)
	uid, _ := strconv.ParseInt(fmt.Sprintf("%v", userInfo["id"]), 10, 32)

	post.UserID = uint(uid)

	if err := repo.Update(&post); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, post)

	sendMessage(s, "post", Message{
		Status: "updated",
		Data:   post,
	})
}

//DeletePost delete single post
func DeletePost(w http.ResponseWriter, r *http.Request, db *gorm.DB, s *sse.Server) {
	postID := chi.URLParam(r, "PostID")
	repo := &repo.PostRepo{Db: db}
	id, _ := strconv.ParseUint(postID, 10, 32)

	if err := repo.DeleteByID(uint(id)); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, postID)

	sendMessage(s, "post", Message{
		Status: "deleted",
		Data:   postID,
	})
}
