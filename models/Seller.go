package models

import (
    "gorm.io/gorm"
)

type (
    Seller struct {
        ID          uint      `gorm:"primary_key" json:"id"`
        Name        string    `json:"name"`
        Address 	string    `json:"address"`
		UserID 		uint `json:"userID"`
        User      	User   `json:"-"`
        Product     []Product   `json:"-"`
        Transaction      	[]Transaction  	 `json:"-"`
    }
)

func ExtractSeller(id uint, db *gorm.DB) (uint, error) {
    // get db from gin context
    var seller Seller
    if err := db.Where("user_id = ?", id).First(&seller).Error; err != nil {
        return 0, err
    }

    idSeller := seller.ID

    return idSeller, nil
}