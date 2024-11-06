package server

import (
	"Rzeczodzielnia/internal/models"
	"Rzeczodzielnia/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterRequest struct {
	Name     string `validate:"required" json:"name"`
	Lastname string `validate:"required" json:"lastname"`
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

type LoginRequest struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	user := models.GetUserByEmail(requestBody.Email)
	if user == nil {
		handleError(w, http.StatusNotFound, "User not found")
		return
	}
	isPassword := utils.ComparePassword(requestBody.Password, user.Password)
	if !isPassword {
		handleError(w, http.StatusUnauthorized, "Invalid password")
		return
	}
	token, err := utils.CreateToken(user.ID, user.Email, user.Role)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error creating token")
		return
	}
	msg := make(map[string]string)
	msg["message"] = "User logged in"
	msg["token"] = token
	sendJSONResponse(w, http.StatusOK, msg)
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}

	if requestBody.Email == models.GetUserByEmail(requestBody.Email).Email {
		handleError(w, http.StatusBadRequest, "User already exists")
		return
	}
	password, err := utils.EncryptPassword(requestBody.Password)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Problems with password encryption")
		return
	}

	newUser := models.User{
		Role:     0,
		Name:     requestBody.Name,
		Lastname: requestBody.Lastname,
		Email:    requestBody.Email,
		Password: password,
		Address:  models.Address{Country: "Default"},
	}
	utils.AddOrUpdateObject(newUser, false)
	msg := make(map[string]string)
	msg["message"] = "User created"
	msg["email"] = newUser.Email
	msg["fullName"] = fmt.Sprintf("%s %s", newUser.Name, newUser.Lastname)

	sendJSONResponse(w, http.StatusCreated, msg)
}
