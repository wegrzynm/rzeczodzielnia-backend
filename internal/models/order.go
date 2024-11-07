package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID         uint    `json:"userID"`
	User           User    `json:"user"`
	CartID         uint    `json:"cartID"`
	Cart           Cart    `json:"cart"`
	PaymentStatus  bool    `json:"paymentStatus"`  // false - not paid, true - paid
	DeliveryStatus uint    `json:"deliveryStatus"` // 0 - not delivered, 1 - package ready, 2- delivery in progress 3- delivered
	OrderAddressId uint    `json:"orderAddressId"`
	OrderAddress   Address `json:"orderAddress"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&Order{}, &Cart{}, &User{}, &Address{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllOrders() []Order {
	var orders []Order
	database.DbInstance.Db.Find(&orders)
	return orders
}

func GetOrderById(Id uint) *Order {
	var order Order
	database.DbInstance.Db.Where("Id=?", Id).First(&order)
	return &order
}
