package repo

import (
	"sc-app/model"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type PostRepo struct {
	Db *gorm.DB
}

//Create new post
func (p *PostRepo) Create(post *model.Post) error {
	//validate struct
	if _, err := govalidator.ValidateStruct(post); err != nil {
		return err
	}
	//reset the ID
	post.ID = 0

	return p.Db.Save(post).Error
}

//FindByID get post by ID
func (p *PostRepo) FindByID(id uint) (*model.Post, error) {
	post := model.Post{}
	rest := p.Db.First(&post, id)
	if rest.RecordNotFound() {
		return nil, &RecordNotFoundError{}
	}

	return &post, rest.Error
}

//Update post
func (p *PostRepo) Update(post *model.Post) error {
	//validate struct
	if _, err := govalidator.ValidateStruct(post); err != nil {
		return err
	}

	foundPost, err := p.FindByID(post.ID)
	if err != nil {
		return err
	}

	post.UserID = foundPost.UserID

	return p.Db.Save(post).Error
}

//DeleteByID remove post by id
func (p *PostRepo) DeleteByID(id uint) error {
	res := p.Db.Delete(model.Post{}, "id = ? ", id)
	if res.RowsAffected == 0 {
		return &RecordNotFoundError{}
	}
	return res.Error
}

//All Get all post order by id DESC
func (p *PostRepo) All() ([]model.Post, error) {
	posts := []model.Post{}
	res := p.Db.Order("id DESC").Find(&posts)
	return posts, res.Error
}
