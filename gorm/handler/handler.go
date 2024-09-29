package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang2/repo"
)

type Handler struct {
	repo *repo.Repository
}

func NewHandler(repo *repo.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user := r.Group("/users")
	{
		user.POST("/", h.createUser)
		user.DELETE("/:id", h.deleteUser)
		user.PUT("/:id", h.updateUser)
		user.GET("/:id", h.getUserById)
		user.GET("/", h.getAllUsers)
	}

	return r
}
