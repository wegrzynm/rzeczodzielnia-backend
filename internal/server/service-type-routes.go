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

func (s *Server) GetServiceTypesHandler(w http.ResponseWriter, r *http.Request) {
	serviceTypes := models.GetAllServiceTypes()

	sendJSONResponse(w, http.StatusOK, &serviceTypes)
}

func (s *Server) CreateServiceTypeHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	if usr.Role != 1 {
		handleError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	var requestBody models.ServiceType
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	if requestBody.Name == models.GetServiceTypeByName(requestBody.Name).Name {
		handleError(w, http.StatusBadRequest, "Service type already exists")
		return
	}
	utils.AddOrUpdateObject(&requestBody, false)
	sendJSONResponse(w, http.StatusCreated, &requestBody)
}

func (s *Server) UpdateServiceTypeHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	if usr.Role != 1 {
		handleError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	serviceTypeId := getParamsId(params)
	serviceType := models.GetServiceTypeById(serviceTypeId)
	if serviceType.ID == 0 {
		handleError(w, http.StatusNotFound, fmt.Sprintf("Service Type not found"))
		return
	}

	var requestBody models.ServiceType
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}

	if requestBody.Name == models.GetServiceTypeByName(requestBody.Name).Name {
		handleError(w, http.StatusBadRequest, "Category already exists")
		return
	}

	serviceType.UpdatedAt = time.Now()
	serviceType.Name = requestBody.Name
	serviceType.Description = requestBody.Description
	utils.AddOrUpdateObject(&serviceType, true)

	sendJSONResponse(w, http.StatusOK, &serviceType)
}
