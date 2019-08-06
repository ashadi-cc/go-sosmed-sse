package repo

import (
	"log"
	"os"
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

func TestMain(m *testing.M) {
	DB, Err = config.ConnectDB()
	if Err != nil {
		log.Fatal(Err)
	}
	defer DB.Close()

	DB.DropTableIfExists(&model.Post{}, &model.User{})

	DB = model.DbMigrate(DB)
	DB.LogMode(false)

	code := m.Run()

	os.Exit(code)
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

	oldID := user.ID
	user.Email = "ashadi2@gmail.com"
	user.Password = "ashadi"
	p.Create(&user)
	if oldID == user.ID {
		t.Fatal("its should not update existing user ", oldID, user.ID)
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

func TestCreatePost(t *testing.T) {
	p := &PostRepo{Db: DB}
	post := model.Post{}
	user := model.User{}
	DB.First(&user)

	post.Title = "bajigan"
	post.Body = "tolong"
	post.UserID = uint(user.ID)

	if err := p.Create(&post); err != nil {
		t.Fatalf("it should created new post but got %s", err.Error())
	}

	oldID := post.ID
	p.Create(&post)
	if oldID == post.ID {
		t.Fatalf("it should not update existing record %d != %d", oldID, post.ID)
	}

	post.UserID = 10
	if err := p.Create(&post); err == nil {
		t.Fatal("it should not create post because the user is not exist")
	}
}

func TestGetAllPOst(t *testing.T) {
	p := &PostRepo{Db: DB}
	posts, err := p.All()
	if err != nil {
		t.Fatalf("got error %s when pull all post", err.Error())
	}

	if len(posts) == 0 {
		t.Fatal("posts should get count > 0")
	}
}

func TestFindPostByID(t *testing.T) {
	p := &PostRepo{Db: DB}

	post, err := p.FindByID(1)
	if err != nil {
		t.Fatalf("it should get post but got error %s", err.Error())
	}

	if post.ID != 1 {
		t.Fatalf("post ID should 1 but got %d", post.ID)
	}
}

func TestUpdatePost(t *testing.T) {
	p := &PostRepo{Db: DB}
	post := model.Post{UserID: 1}
	post.ID, post.Title, post.Body = 1, "abc", "def"

	if err := p.Update(&post); err != nil {
		t.Fatalf("it should update post but got error %s", err.Error())
	}

	if post.Title != "abc" {
		t.Fatalf("the title should be abc but got %s", post.Title)
	}

	post.ID = 20
	if err := p.Update(&post); err == nil {
		t.Fatalf("its should not created new post %v", post)
	}
}

func TestDeletePost(t *testing.T) {
	p := &PostRepo{Db: DB}
	if err := p.DeleteByID(1); err != nil {
		t.Fatalf("it should delete post but got error %s", err.Error())
	}

	post, _ := p.FindByID(1)
	if post != nil {
		t.Fatalf("the post with id 1 should delete but still exist %v", post)
	}

	if err := p.DeleteByID(10); err == nil {
		t.Fatalf("it should got error but got error %s", err.Error())
	}
}
