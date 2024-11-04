package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllCompanies() []User {
	var company []User
	database.DbInstance.Db.Find(&company)
	return company
}

func GetUserById(Id uint) *User {
	var user User
	database.DbInstance.Db.Where("Id=?", Id).Find(&user)
	return &user
}

func GetUserByEmail(email string) *User {
	var user User
	database.DbInstance.Db.Where("email=?", email).First(&user)
	return &user
}
