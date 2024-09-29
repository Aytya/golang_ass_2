package main

import (
	"golang2/database"
	_ "golang2/docs"
	"golang2/handler"
	"golang2/repo"
	"log"
)

// @title           Service for users, profile
// @version         1.0
// @description     RESTful API

// @host      localhost:8080
// @BasePath  /
func main() {
	db := database.Init()
	repository := repo.NewRepository(db)
	handler := handler.NewHandler(repository)

	routes := handler.InitRoutes()
	if err := routes.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
