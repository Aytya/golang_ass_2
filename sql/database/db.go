package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var (
	usersTable = "users"
)

var db *sql.DB

func Init() *sql.DB {
	dbURL := "user=postgres password=7212Hey) dbname=postgres sslmode=disable"

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	if err := createUserTable(); err != nil {
		log.Fatalln(err)
	}

	return db
}

func createUserTable() error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
				ID SERIAL PRIMARY KEY, 
					name VARCHAR(255) NOT NULL UNIQUE,
					age INT NOT NULL
			)`, usersTable)

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
