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

type ProductResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	CategoryID  uint   `json:"categoryID"`
}

func (s *Server) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	products := models.GetAllProducts(true)

	sendJSONResponse(w, http.StatusOK, products)
}

func (s *Server) GetProductsByCategoryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	categoryId := getParamsId(params)
	products := models.GetProductsByCategory(categoryId)

	sendJSONResponse(w, http.StatusOK, products)
}

func (s *Server) GetProductsByUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := getParamsId(params)
	products := models.GetProductsByUser(userId)

	sendJSONResponse(w, http.StatusOK, products)
}

func (s *Server) GetProductHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	productId := getParamsId(params)
	product := models.GetProductById(productId)

	if product.ID == 0 {
		handleError(w, http.StatusNotFound, fmt.Sprintf("Product not found"))
		return
	}

	sendJSONResponse(w, http.StatusOK, product)
}

func (s *Server) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	var requestBody models.Product
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	requestBody.UserID = usr.Id

	utils.AddOrUpdateObject(&requestBody, false)
	sendJSONResponse(w, http.StatusCreated, requestBody)
}

func (s *Server) UpdateProductHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	productId := getParamsId(params)
	product := models.GetProductById(productId)
	if product.ID == 0 {
		handleError(w, http.StatusNotFound, fmt.Sprintf("Product not found"))
		return
	}
	if product.UserID != usr.Id {
		handleError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	var requestBody models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}

	if requestBody.Name != "" {
		product.Name = requestBody.Name
	}
	if requestBody.Description != "" {
		product.Description = requestBody.Description
	}
	if requestBody.CategoryID != 0 {
		product.CategoryID = requestBody.CategoryID
	}
	product.UpdatedAt = time.Now()
	utils.AddOrUpdateObject(&product, true)
	sendJSONResponse(w, http.StatusOK, product)
}

func (s *Server) DeleteProductHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	productId := getParamsId(params)
	product := models.GetProductById(productId)
	if product.ID == 0 {
		handleError(w, http.StatusNotFound, fmt.Sprintf("Product not found"))
		return
	}
	if product.UserID != usr.Id {
		handleError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	product.IsActive = false
	utils.AddOrUpdateObject(&product, true)
	models.DeleteProductById(*product)
	msg := map[string]string{"message": "Product deleted successfully"}
	sendJSONResponse(w, http.StatusNoContent, msg)
}
