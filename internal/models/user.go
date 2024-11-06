package models

import (
	"Rzeczodzielnia/internal/database"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string  `gorm:"unique" json:"email"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Lastname  string  `json:"lastname"`
	AddressId uint    `json:"addressId"`
	Address   Address `json:"address"`
	Role      uint    `json:"role"` // 0 - user, 1 - admin
}

func init() {
	database.New()
	err := database.DbInstance.Db.AutoMigrate(&Address{}, &User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetAllUsers() []User {
	var users []User
	database.DbInstance.Db.Where("role=?", 0).Find(&users)
	return users
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
