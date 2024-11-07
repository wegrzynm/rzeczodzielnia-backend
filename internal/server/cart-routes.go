package server

import (
	"Rzeczodzielnia/internal/models"
	"Rzeczodzielnia/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"math"
	"net/http"
	"time"
)

type CartItems struct {
	Products []uint `json:"products"`
}

type PromoCodeResponse struct {
	PromoCode string `json:"promoCode"`
}

func (s *Server) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	cart := *models.GetCartByUserId(usr.Id)
	if cart.ID == 0 {
		cart = models.Cart{
			UserID:       usr.Id,
			IsCheckedOut: false,
		}
		utils.AddOrUpdateObject(&cart, false)
	}

	sendJSONResponse(w, http.StatusOK, &cart)
}

func (s *Server) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	cart := models.GetCartByUserId(usr.Id)
	if cart.ID == 0 {
		cart = &models.Cart{
			UserID:       usr.Id,
			IsCheckedOut: false,
		}
		utils.AddOrUpdateObject(&cart, false)
	}

	var requestBody CartItems
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	for _, v := range requestBody.Products {
		product := models.GetActiveProductById(v, true)
		if product.ID == 0 {
			handleError(w, http.StatusBadRequest, fmt.Sprintf("Product with id %d not found", v))
			return
		}
		cart.Total += product.Price
		cart.Items = append(cart.Items, *product)
	}

	cart.UpdatedAt = time.Now()

	utils.AddOrUpdateObject(&cart, true)
	sendJSONResponse(w, http.StatusOK, &cart)
}

func (s *Server) RemoveFromCartHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	productId := getParamsId(params)
	cart := models.GetCartByUserId(usr.Id)
	if cart.ID == 0 {
		cart = &models.Cart{
			UserID:       usr.Id,
			IsCheckedOut: false,
		}
		utils.AddOrUpdateObject(&cart, false)
	}

	for i, item := range cart.Items {
		if item.ID == productId {
			cart.Total -= item.Price
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}
	models.RemoveItemFromCart(cart)
	cart.UpdatedAt = time.Now()
	utils.AddOrUpdateObject(&cart, true)
	sendJSONResponse(w, http.StatusOK, &cart)
}

func (s *Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	cart := models.GetCartByUserId(usr.Id)
	if cart.ID == 0 {
		cart = &models.Cart{
			UserID:       usr.Id,
			IsCheckedOut: false,
		}
		utils.AddOrUpdateObject(&cart, false)
	}

	cart.Items = nil
	cart.Total = 0
	cart.UpdatedAt = time.Now()
	models.RemoveItemFromCart(cart)
	utils.AddOrUpdateObject(&cart, true)
	sendJSONResponse(w, http.StatusOK, &cart)
}

func (s *Server) AddPromoCodeCartHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	cart := models.GetCartByUserId(usr.Id)
	if cart.ID == 0 {
		handleError(w, http.StatusBadRequest, "Cart not found")
		return
	}
	var requestBody PromoCodeResponse
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	promoCode := models.GetPromoCodeByCode(requestBody.PromoCode)
	if promoCode.ID == 0 {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Promo code %s not found", requestBody.PromoCode))
		return
	}
	if cart.PromoCode != "" {
		handleError(w, http.StatusBadRequest, "Promo code already applied")
		return
	}
	cart.PromoCode = promoCode.Code
	cart.Total = math.Round(cart.Total*(1-(promoCode.Discount/100))*100) / 100
	cart.UpdatedAt = time.Now()
	utils.AddOrUpdateObject(&cart, true)
	sendJSONResponse(w, http.StatusOK, &cart)
}

func (s *Server) RemovePromoCodeCartHandler(w http.ResponseWriter, r *http.Request) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	cart := models.GetCartByUserId(usr.Id)
	if cart.ID == 0 {
		handleError(w, http.StatusBadRequest, "Cart not found")
		return
	}
	var requestBody PromoCodeResponse
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err))
		return
	}
	promoCode := models.GetPromoCodeByCode(requestBody.PromoCode)
	if promoCode.ID == 0 {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Promo code %s not found", requestBody.PromoCode))
		return
	}
	cart.PromoCode = ""
	cart.Total = math.Round(cart.Total/(1-(promoCode.Discount/100))*100) / 100
	cart.UpdatedAt = time.Now()
	utils.AddOrUpdateObject(&cart, true)
	sendJSONResponse(w, http.StatusOK, &cart)
}
