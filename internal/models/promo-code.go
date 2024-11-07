package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type PromoCode struct {
	gorm.Model
	Code     string  `json:"code"`
	Discount float64 `json:"discount"` // in percentage
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&PromoCode{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllPromoCodes() []PromoCode {
	var promoCodes []PromoCode
	database.DbInstance.Db.Find(&promoCodes)
	return promoCodes
}

func GetPromoCodeById(Id uint) *PromoCode {
	var promoCode PromoCode
	database.DbInstance.Db.Where("id = ?", Id).First(&promoCode)
	return &promoCode
}

func GetPromoCodeByCode(code string) *PromoCode {
	var promoCode PromoCode
	database.DbInstance.Db.Where("code = ?", code).First(&promoCode)
	return &promoCode
}
