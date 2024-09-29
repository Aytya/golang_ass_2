package repo

import (
	"database/sql"
	"sql/model"
)

type User interface {
	InsertUser(users []model.User) error
	FindUsers(age int, limit, offset int) ([]model.User, error)
	UpdateUserById(newName string, newAge int, id uint) error
	DeleteUserById(id uint) error
}

type Repository struct {
	User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
