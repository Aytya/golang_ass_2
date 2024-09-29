package repo

import (
	"database/sql"
	"fmt"
	"sql/model"
)

type UserRepository struct {
	db *sql.DB
}

var (
	usersTable = "users"
	usersName  = "name"
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) InsertUser(users []model.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.Prepare("INSERT INTO users(name, age) VALUES($1, $2)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, user := range users {
		var exists bool
		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", user.Name).Scan(&exists)
		if err != nil {
			tx.Rollback()
			return err
		}
		if exists {
			tx.Rollback()
			return fmt.Errorf("user with name '%s' already exists", user.Name)
		}

		_, err = query.Exec(user.Name, user.Age)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindUsers(age, limit, offset int) ([]model.User, error) {
	var users []model.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE age = %d ORDER BY %s ASC LIMIT %d OFFSET %d", usersTable, age, usersName, limit, offset)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over users: %v", err)
	}

	return users, nil
}

func (r *UserRepository) UpdateUserById(newName string, newAge int, id uint) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = $1 AND id != $2)", newName, id).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user with name '%s' already exists", newName)
	}

	query := fmt.Sprintf("UPDATE users SET name = $1, age = $2 WHERE id = $3")
	_, err = r.db.Exec(query, newName, newAge, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUserById(id uint) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
