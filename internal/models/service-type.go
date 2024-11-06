package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type ServiceType struct {
	gorm.Model
	Name        string `gorm:"unique" json:"name"`
	Description string `json:"description"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&ServiceType{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllServiceTypes() []ServiceType {
	var serviceTypes []ServiceType
	database.DbInstance.Db.Find(&serviceTypes)
	return serviceTypes
}

func GetServiceTypeById(Id uint) *ServiceType {
	var serviceType ServiceType
	database.DbInstance.Db.Where("Id=?", Id).Find(&serviceType)
	return &serviceType
}

func GetServiceTypeByName(name string) *ServiceType {
	var serviceType ServiceType
	database.DbInstance.Db.Where("name=?", name).Find(&serviceType)
	return &serviceType
}
