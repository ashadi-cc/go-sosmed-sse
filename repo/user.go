package repo

import (
	"sc-app/model"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//UserRepo repository
type UserRepo struct {
	Db *gorm.DB
}

//Create New User
func (u *UserRepo) Create(user *model.User) error {
	//validate struct
	if _, err := govalidator.ValidateStruct(user); err != nil {
		return err
	}
	//check email exist or not
	if _, err := u.FindByEmail(user.Email); err == nil {
		return &EmailExistError{}
	}

	//reset the ID
	user.ID = 0
	//hash password
	user.Password, _ = GeneratePassword(user.Password)
	err := u.Db.Save(&user).Error
	return err
}

//FindByEmail find user by email
func (u *UserRepo) FindByEmail(email string) (*model.User, error) {
	user := model.User{}
	res := u.Db.Where("email = ?", email).First(&user)

	if res.RecordNotFound() {
		return nil, &EmailNotExistsError{}
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

//Login user
func (u *UserRepo) Login(user *model.User) error {
	//validate struct
	if _, err := govalidator.ValidateStruct(user); err != nil {
		return err
	}

	loginUser, err := u.FindByEmail(user.Email)
	if err != nil {
		return err
	}

	if !MatchPassword(loginUser.Password, user.Password) {
		return &PasswordError{}
	}

	user.ID = loginUser.ID
	user.Password = ""
	return nil
}

//GeneratePassword with bycript function
func GeneratePassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(passwordHash), err
}

//MatchPassword //compare password
func MatchPassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}
