package utils

import (
	"Rzeczodzielnia/internal/database"
	"log"
)

func AddOrUpdateObject[T any](object T, isUpdate bool) {
	if isUpdate {
		result := database.DbInstance.Db.Save(&object)
		if result.Error != nil {
			log.Printf("Error updating object: %s", result.Error.Error())
			return
		}
		log.Printf("Updated object")
		return
	}

	result := database.DbInstance.Db.Create(&object)
	if result.Error != nil {
		log.Printf("Error adding object: %s", result.Error.Error())
		return
	}
	log.Printf("Added object")
}
