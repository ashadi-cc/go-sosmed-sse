package model

type Post struct {
	CommonModel
	UserId uint `json:"-"`
	Title string `json:"title"`
	Body string `json:"body" gorm:"type:text"`
}