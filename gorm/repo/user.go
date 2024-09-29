package repo

import (
	"golang2/model"
	"gorm.io/gorm"
)

type UserPostgres struct {
	db *gorm.DB
}

func (u *UserPostgres) CreateUser(user model.User) (string, error) {
	if err := u.db.Create(&user).Error; err != nil {
		return "", err
	}

	return "User successfully created", nil
}

func (u *UserPostgres) UpdateUser(id string, user model.User) error {
	if err := u.db.Model(&user).Where("id = ?", id).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) DeleteUser(id string) error {
	if err := u.db.Delete(&model.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) GetUserById(id string) (model.User, error) {
	var user model.User
	if err := u.db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err

	}
	return user, nil
}

func (u *UserPostgres) GetAllUsers(age int, limit, offset int) ([]model.User, error) {
	var users []model.User
	query := u.db.Model(&model.User{}).Where("age = ?", age).Order("name desc")

	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}
