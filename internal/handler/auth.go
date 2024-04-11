package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		response.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	jsonResponse := response.UserResponse201{ID: id}
	err = json.NewEncoder(w).Encode(jsonResponse)
	if err != nil {
		log.Println(err)
	}
}

type signInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.Role != "admin" && input.Role != "user" {
		response.HandleError(w, &response.AccessError{Message: "incorrect role"})
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password, input.Role)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	jsonResponse := map[string]interface{}{
		"token": token,
	}
	err = json.NewEncoder(w).Encode(jsonResponse)
	if err != nil {
		log.Println(err)
	}
}
