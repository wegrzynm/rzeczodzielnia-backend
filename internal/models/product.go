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
	Price       float64  `json:"price"`
	Quantity    int      `json:"quantity"`
	CategoryID  uint     `json:"categoryID"`
	Category    Category `json:"category"`
	Images      []Image  `json:"images" gorm:"foreignKey:ProductID"`
	UserID      uint     `json:"userID"`
	User        User     `json:"user"`
	IsActive    bool     `json:"isActive"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&Product{}, &Category{}, &Image{}, &User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllProducts(isActive bool) []Product {
	var products []Product
	database.DbInstance.Db.Preload("Category").Preload("Images").
		Where("is_active=?", isActive).
		Find(&products)
	return products
}

func GetProductsByCategory(category uint) []Product {
	var products []Product
	database.DbInstance.Db.Preload("Category").Preload("Images").
		Where("category_id=?", category).Where("is_active=?", true).
		Find(&products)
	return products
}

func GetProductsByUser(userId uint) []Product {
	var products []Product
	database.DbInstance.Db.Preload("Category").Preload("Images").
		Where("user_id=?", userId).Where("is_active=?", true).
		Find(&products)
	return products
}

func GetProductById(Id uint) *Product {
	var product Product
	database.DbInstance.Db.Preload("Category").Preload("Images").Where("Id=?", Id).First(&product)
	return &product
}

func DeleteProductById(product Product) {
	database.DbInstance.Db.Delete(&product)
}
