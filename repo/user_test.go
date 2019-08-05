package repo

import (
	"sc-app/config"
	"sc-app/model"
	"testing"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DB  *gorm.DB
	Err error
)

func TestConnection(t *testing.T) {
	DB, Err = config.ConnectDB()
	if Err != nil {
		t.Fatalf("cannot connect to database %s", Err.Error())
	}

	DB.DropTableIfExists(&model.User{})

	DB = model.DbMigrate(DB)
}

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword("123456")

	if err != nil {
		t.Fatalf("Error when generating password: %s", err.Error())
	}

	if !MatchPassword(password, "123456") {
		t.Fatalf("Password does not match")
	}

}

func TestValidatorUser(t *testing.T) {
	p := &UserRepo{Db: DB}

	user := model.User{}

	if err := p.Create(&user); err == nil {
		t.Fatal("its should error because email and password empty")
	}

	user.Email = "aa"
	user.Password = ""

	if err := p.Create(&user); err == nil {
		t.Fatal("its should error because email is not valid and password empty")
	}

	user.Email = "ashadi.cc@gmail.com"
	user.Password = ""

	if err := p.Create(&user); err == nil {
		t.Fatal("its should error because password is empty")
	}

	user.Email = ""
	user.Password = "12342323"

	if err := p.Create(&user); err == nil {
		t.Fatal("its should error because email is empty")
	}

	user.Email = "aaa"
	user.Password = "121212"

	if err := p.Create(&user); err == nil {
		t.Fatal("its should error because email is invalid")
	}

}
func TestShouldCreateUser(t *testing.T) {
	p := &UserRepo{Db: DB}

	DB.Exec("DELETE FROM users WHERE email = ? ", "ashadi.cc@gmail.com")

	user := model.User{}
	user.Email = "ashadi.cc@gmail.com"
	user.Password = "123456"

	err := p.Create(&user)

	if err != nil {
		t.Fatalf("Cannot create user get error %s", err.Error())
	}

	if user.Password == "123456" {
		t.Fatalf("password does not hashed, %s", user.Password)
	}
}

func TestFindUserByEmail(t *testing.T) {
	p := &UserRepo{Db: DB}

	user, err := p.FindByEmail("ashadi.cc@gmail.com")

	if err != nil {
		t.Fatalf("cannot find user with email ashadi.cc@gmail.com: %s", err.Error())
	}

	if user.Email != "ashadi.cc@gmail.com" {
		t.Fatalf("Expected email ashadi.cc@gmail.com but got %s", user.Email)
	}
}

func TestLoginUser(t *testing.T) {
	p := &UserRepo{Db: DB}
	user := &model.User{}
	user.Email = "ashadi.cc@gmail.com"
	user.Password = "123456"

	if err := p.Login(user); err != nil {
		t.Fatalf("login error: %s", err.Error())
	}

}
