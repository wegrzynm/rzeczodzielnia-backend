package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	ProductID uint   `json:"productID"`
	Name      string `gorm:"unique" json:"name"`
	Path      string `json:"path"`
	UserId    uint   `json:"userId"`
	User      User   `json:"user"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&Image{}, &Product{}, &User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllImages() []Image {
	var images []Image
	database.DbInstance.Db.Find(&images)
	return images
}

func GetImageById(Id uint) *Image {
	var image Image
	database.DbInstance.Db.Where("Id=?", Id).First(&image)
	return &image
}

func GetImageByName(name string) *Image {
	var image Image
	database.DbInstance.Db.Where("name=?", name).First(&image)
	return &image
}

func GetImagesByProductId(Id uint) []Image {
	var images []Image
	database.DbInstance.Db.Where("product_id=?", Id).Find(&images)
	return images
}

func DeleteImageById(image Image) {
	database.DbInstance.Db.Delete(&image)
}
