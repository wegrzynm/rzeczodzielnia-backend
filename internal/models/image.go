package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	ProductID   uint   `json:"productID"`
	Name        string `gorm:"unique" json:"name"`
	Path        string `json:"path"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&Image{})
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
	database.DbInstance.Db.Where("Id=?", Id).Find(&image)
	return &image
}
