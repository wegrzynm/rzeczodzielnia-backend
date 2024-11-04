package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Country string `json:"country"`
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
	Street  string `json:"street"`
	Number  string `json:"number"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&Address{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllAddresses() []Address {
	var addresses []Address
	database.DbInstance.Db.Find(&addresses)
	return addresses
}

func GetAddressById(Id uint) *Address {
	var address Address
	database.DbInstance.Db.Where("Id=?", Id).Find(&address)
	return &address
}
