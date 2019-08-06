package model

type Post struct {
	CommonModel
	User   User   `gorm:"association_foreignkey:UserID" valid:"-" json:"-"`
	Title  string `json:"title" valid:"required"`
	Body   string `json:"body" gorm:"type:text" valid:"required"`
	UserID uint   `json:"-"`
}
