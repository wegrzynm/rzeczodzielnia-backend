package server

import (
	"Rzeczodzielnia/internal/models"
	"Rzeczodzielnia/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangeUserDataRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func (s *Server) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	user := models.GetUserById(usr.Id)
	sendJSONResponse(w, http.StatusOK, &user)
}

func (s *Server) SetUserAddressHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	user := models.GetUserById(usr.Id)
	var requestBody models.Address
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}

	if requestBody.Country != "" {
		user.Address.Country = requestBody.Country
	}
	if requestBody.City != "" {
		user.Address.City = requestBody.City
	}
	if requestBody.Street != "" {
		user.Address.Street = requestBody.Street
	}
	if requestBody.ZipCode != "" {
		user.Address.ZipCode = requestBody.ZipCode
	}
	if requestBody.Number != "" {
		user.Address.Number = requestBody.Number
	}

	utils.AddOrUpdateObject(user.Address, true)
	sendJSONResponse(w, http.StatusOK, &user)
}

func (s *Server) ChangeUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	user := models.GetUserById(usr.Id)
	var requestBody ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	isOldPasswordValid := utils.ComparePassword(requestBody.OldPassword, user.Password)
	if !isOldPasswordValid {
		handleError(w, http.StatusBadRequest, "Old password is incorrect")
		return
	}

	encryptedPassword, err := utils.EncryptPassword(requestBody.NewPassword)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Sprintf("Error encrypting password: %v", err))
		return
	}
	user.Password = encryptedPassword
	user.UpdatedAt = time.Now()
	utils.AddOrUpdateObject(user, true)
	msg := map[string]string{"message": "Password changed successfully"}
	sendJSONResponse(w, http.StatusOK, msg)
}

func (s *Server) ChangeUserDataHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	user := models.GetUserById(usr.Id)
	var requestBody ChangeUserDataRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}

	if requestBody.Email != "" {
		user.Email = requestBody.Email
	}
	if requestBody.Name != "" {
		user.Name = requestBody.Name
	}
	if requestBody.Lastname != "" {
		user.Lastname = requestBody.Lastname
	}

	user.UpdatedAt = time.Now()
	utils.AddOrUpdateObject(user, true)
	sendJSONResponse(w, http.StatusOK, &user)
}
