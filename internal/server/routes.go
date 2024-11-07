package server

import (
	"Rzeczodzielnia/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()
	routeVersion := "/v1"
	//Login and Register routes
	r.HandlerFunc(http.MethodPost, routeVersion+"/login", s.LoginHandler)
	r.HandlerFunc(http.MethodPost, routeVersion+"/register", s.RegisterHandler)

	//Category routes
	r.HandlerFunc(http.MethodGet, routeVersion+"/categories", s.GetCategoriesHandler)
	r.HandlerFunc(http.MethodPost, routeVersion+"/category/create", s.CreateCategoryHandler)
	r.Handle(http.MethodPut, routeVersion+"/category/:id/update", s.UpdateCategoryHandler)

	//Service routes
	r.HandlerFunc(http.MethodGet, routeVersion+"/service-types", s.GetServiceTypesHandler)
	r.HandlerFunc(http.MethodPost, routeVersion+"/service-type/create", s.CreateServiceTypeHandler)
	r.Handle(http.MethodPut, routeVersion+"/service-type/:id/update", s.UpdateServiceTypeHandler)

	//Product routes
	r.HandlerFunc(http.MethodGet, routeVersion+"/products", s.GetProductsHandler)
	r.Handle(http.MethodGet, routeVersion+"/products/category/:id", s.GetProductsByCategoryHandler)
	r.Handle(http.MethodGet, routeVersion+"/products/user/:id", s.GetProductsByUserHandler)
	r.Handle(http.MethodGet, routeVersion+"/product/:id", s.GetProductHandler)
	r.HandlerFunc(http.MethodPost, routeVersion+"/product/create", s.CreateProductHandler)
	r.Handle(http.MethodPut, routeVersion+"/product/:id/update", s.UpdateProductHandler)
	r.Handle(http.MethodDelete, routeVersion+"/product/:id/delete", s.DeleteProductHandler)

	//Image routes
	r.Handle(http.MethodPost, routeVersion+"/images/product/:id", s.UploadImageHandler)
	r.Handle(http.MethodDelete, routeVersion+"/images/:id/delete", s.DeleteImageHandler)

	r.HandlerFunc(http.MethodGet, "/", s.HelloWorldHandler)

	r.HandlerFunc(http.MethodGet, "/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, _ := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, _ := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}

func handleError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errMsg := map[string]string{"error": errorMessage}
	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResp, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResp)
}

func userAuthentication(tokenString string) (bool, int, string, *utils.TokenClaims) {
	if tokenString == "" {
		return false, http.StatusUnauthorized, "Missing authorization header", nil
	}
	tokenString = tokenString[len("Bearer "):]

	err, usr := utils.VerifyToken(tokenString)
	if err != nil {
		return false, http.StatusUnauthorized, "Invalid token", nil
	}
	return true, http.StatusOK, "", usr
}

func getParamsId(params httprouter.Params) uint {
	idString := params.ByName("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0
	}
	return uint(id)
}
