package server

import (
	"Rzeczodzielnia/internal/models"
	"Rzeczodzielnia/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (s *Server) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories := models.GetAllCategories()

	sendJSONResponse(w, http.StatusOK, &categories)
}

func (s *Server) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	if usr.Role != 1 {
		handleError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	var requestBody models.Category
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	if requestBody.Name == models.GetCategoryByName(requestBody.Name).Name {
		handleError(w, http.StatusBadRequest, "Category already exists")
		return
	}
	utils.AddOrUpdateObject(&requestBody, false)
	sendJSONResponse(w, http.StatusCreated, &requestBody)
}

func (s *Server) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	if usr.Role != 1 {
		handleError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	categoryId := getParamsId(params)
	category := models.GetCategoryId(categoryId)
	if category.ID == 0 {
		handleError(w, http.StatusNotFound, fmt.Sprintf("Category not found"))
		return
	}

	var requestBody models.Category
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}

	if requestBody.Name == models.GetCategoryByName(requestBody.Name).Name {
		handleError(w, http.StatusBadRequest, "Category already exists")
		return
	}

	category.UpdatedAt = time.Now()
	category.Name = requestBody.Name
	category.Description = requestBody.Description
	utils.AddOrUpdateObject(&category, true)
	sendJSONResponse(w, http.StatusOK, &category)
}
