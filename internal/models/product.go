package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Quantity    int      `json:"quantity"`
	CategoryID  uint     `json:"categoryID"`
	Category    Category `json:"category"`
	Images      []Image  `json:"images" gorm:"foreignKey:ProductID"`
	UserID      uint     `json:"userID"`
	User        User     `json:"user"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&Product{}, &Category{}, &Image{}, &User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllProducts() []Product {
	var products []Product
	database.DbInstance.Db.Preload("Category").Preload("Images").Find(&products)
	return products
}

func GetProductsByCategory(category Category) []Product {
	var products []Product
	database.DbInstance.Db.Preload("Category").Preload("Images").Where("category_id=?", category.ID).Find(&products)
	return products
}

func GetProductsByUser(user User) []Product {
	var products []Product
	database.DbInstance.Db.Preload("Category").Preload("Images").Where("user_id=?", user.ID).Find(&products)
	return products
}

func GetProductById(Id uint) *Product {
	var product Product
	database.DbInstance.Db.Preload("Category").Preload("Images").Where("Id=?", Id).Find(&product)
	return &product
}
