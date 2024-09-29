package repo

import (
	"golang2/model"
	"gorm.io/gorm"
)

type User interface {
	CreateUser(user model.User) (string, error)
	UpdateUser(id string, user model.User) error
	DeleteUser(id string) error
	GetUserById(id string) (model.User, error)
	GetAllUsers(age int, limit, offset int) ([]model.User, error)
}

type Profile interface {
	CreateUserAndProfile(user model.User, profile model.Profile) error
	GetUsersWithProfile() ([]model.Profile, error)
	UpdateUserProfile(userID string, bio string, profilePictureURL string) error
	DeleteUserWithProfile(userID string) error
}

type Repository struct {
	User
	Profile
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:    NewUserPostgres(db),
		Profile: NewProfilePostgres(db),
	}
}
