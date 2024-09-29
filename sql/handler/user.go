package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sql/model"
	"strconv"
)

// InsertUsers InsertUser
// @Summary      add newUser
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body []model.User true "list of users"
// @Success      200  {array}  model.User
// @Router       /users [post]
func (h *Handler) InsertUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
	}

	var users []model.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = h.repo.InsertUser(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// FindUsers
// @Summary      get all users
// @Tags         users
// @Produce      json
// @Param 	     age query string false "filtering by age"
// @Param        limit query string false "pagination per page limit"
// @Param		 offset query string false "pagination offset"
// @Success      200  {array}  model.User
// @Router       /users [get]
func (h *Handler) FindUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	ageStr := r.URL.Query().Get("age")
	age, _ := strconv.Atoi(ageStr)

	var err error
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	users, err := h.repo.FindUsers(age, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// UpdateUserById
// @Summary      update user by id
// @Tags         users
// @Param 		 id path string true "user's id"
// @Param		 name query string true "user's name"
// @Param		 age query string true "user's age"
// @Success      204  "No Content"
// @Router       /users/{id} [put]
func (h *Handler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	newName := r.URL.Query().Get("name")
	if newName == "" {
		http.Error(w, "Missing user name", http.StatusBadRequest)
		return
	}

	ageStr := r.URL.Query().Get("age")
	newAge, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateUserById(newName, newAge, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUserById
// @Summary      delete user by id
// @Tags         users
// @Param 		 id path string false "user's id"
// @Success      204 "No Content"
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteUserById(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
