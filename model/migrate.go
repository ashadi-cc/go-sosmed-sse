package model

import (
	"github.com/jinzhu/gorm"
)

//DbMigrate db migration
func DbMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{}, &Post{})
	db.Model(&Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	return db
}
