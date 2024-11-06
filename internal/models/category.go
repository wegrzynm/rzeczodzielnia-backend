package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string `gorm:"unique" json:"name"`
	Description string `json:"description"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllCategories() []Category {
	var category []Category
	database.DbInstance.Db.Find(&category)
	return category
}

func GetCategoryId(Id uint) *Category {
	var category Category
	database.DbInstance.Db.Where("Id=?", Id).Find(&category)
	return &category
}

func GetCategoryByName(name string) *Category {
	var category Category
	database.DbInstance.Db.Where("name=?", name).Find(&category)
	return &category
}
