package main

import (
	"log"
	"net/http"
	"sql/database"
	_ "sql/docs"
	"sql/handler"
	repo2 "sql/repo"
)

// @title           Service for users
// @version         1.0
// @description     RESTful API

// @host      localhost:8080
// @BasePath  /
func main() {
	db := database.Init()
	repo := repo2.NewRepository(db)
	handler := handler.NewHandler(repo)
	log.Println("Connected to PostgreSQL")
	log.Fatal(http.ListenAndServe(":8080", handler.InitRoutes()))
}
