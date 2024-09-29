package main

import (
	"database/sql"
	"log"
	"sql/model"
	repo2 "sql/repo"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	dsn := "user=testuser dbname=postgres password=testpassword host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT UNIQUE, age INT)")
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}

	return db
}

func tearDownTestDB(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		log.Fatalf("Could not drop users table: %v", err)
	}
}

func TestInsertUser(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)

	repo := repo2.NewRepository(db)

	users := []model.User{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
	}

	err := repo.InsertUser(users)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, len(users), count)

	err = repo.InsertUser([]model.User{{Name: "Alice", Age: 27}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user with name 'Alice' already exists")
}

func TestFindUsers(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)

	repo := repo2.NewRepository(db)

	users := []model.User{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 30},
	}
	_ = repo.InsertUser(users)

	foundUsers, err := repo.FindUsers(30, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, foundUsers, 3)

	foundUsers, err = repo.FindUsers(25, 1, 0)
	assert.NoError(t, err)
	assert.Len(t, foundUsers, 1)
	assert.Equal(t, "Alice", foundUsers[0].Name)
}

func TestUpdateUserById(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)

	repo := repo2.NewRepository(db)

	_, err := db.Exec("INSERT INTO users(name, age) VALUES($1, $2)", "AZA", 25)
	assert.NoError(t, err)

	err = repo.UpdateUserById("Alice Smith", 1, 26)
	assert.NoError(t, err)

	var name string
	var age int
	err = db.QueryRow("SELECT name, age FROM users WHERE id = $1", 1).Scan(&name, &age)
	assert.NoError(t, err)
	assert.Equal(t, "Alice Smith", name)
	assert.Equal(t, 26, age)

	_, _ = db.Exec("INSERT INTO users(name, age) VALUES($1, $2)", "Bob", 30)
	err = repo.UpdateUserById("Bob", 1, 29)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user with name 'Bob' already exists")
}

func TestDeleteUserById(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)

	repo := repo2.NewRepository(db)

	_, err := db.Exec("INSERT INTO users(name, age) VALUES($1, $2)", "Alice", 25)
	assert.NoError(t, err)

	err = repo.DeleteUserById(1)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
