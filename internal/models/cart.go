package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID       uint      `json:"userId"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	Items        []Product `gorm:"many2many:cart_items;" json:"items"`
	PromoCode    string    `json:"promoCode"`
	Total        float64   `json:"total"`
	IsCheckedOut bool      `json:"isCheckedOut"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&User{}, &Product{}, &Cart{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllCarts() []Cart {
	var carts []Cart
	database.DbInstance.Db.Preload("User").Preload("Items").Find(&carts)
	return carts
}

func GetCartById(Id uint) *Cart {
	var cart Cart
	database.DbInstance.Db.Preload("User").Preload("Items").
		Where("id = ?", Id).First(&cart)
	return &cart
}

func GetCartByUserId(Id uint) *Cart {
	var cart Cart
	database.DbInstance.Db.Preload("User").Preload("Items").
		Where("user_id = ?", Id).First(&cart)
	return &cart
}

func RemoveItemFromCart(cart *Cart) {
	database.DbInstance.Db.Model(&cart).Association("Items").Replace(cart.Items)
}
