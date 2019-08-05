package model

import (
	"github.com/jinzhu/gorm"
)

//DbMigrate db migration
func DbMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{}, &Post{})
	return db
}
