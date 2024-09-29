package handler

import (
	"github.com/gin-gonic/gin"
	"golang2/model"
	"net/http"
	"strconv"
)

// AddUser
// @Summary      add newUser
// @Tags         users
// @Accept       json
// @Produce      json
// @Param		 input body model.User true "user request"
// @Success      200  {object}  model.User
// @Router       /users [post]
func (h *Handler) createUser(c *gin.Context) {
	var newUser model.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
		return
	}

	msg, err := h.repo.User.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// GetAllUser
// @Summary      get all users
// @Tags         users
// @Produce      json
// @Param 		 age query string false "filtering by age"
// @Param        limit query string false "pagination per page limit"
// @Param		 offset query string false "pagination offset"
// @Success      200  {array}  model.User
// @Router       /users [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	ageStr := c.Query("age")
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	users, err := h.repo.User.GetAllUsers(age, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteUserByID
// @Summary      delete user by id
// @Tags         users
// @Param 		 id path string true "user's id"
// @Success      204 "OK"
// @Router       /users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if err := h.repo.User.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUserByID
// @Summary      get user by id
// @Tags         users
// @Param 		 id path string true "user's id"
// @Success      204 "OK"
// @Router       /users/{id} [get]
func (h *Handler) getUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := h.repo.User.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateTodo
// @Summary      update user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 id path string true "user's id"
// @Param		 input body model.User true "user's request"
// @Success      204  "No Content"
// @Router       /users/{id} [put]
func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var input model.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
		return
	}

	if err := h.repo.User.UpdateUser(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
