package models

import (
    "gorm.io/gorm"
)

type (
    Customer struct {
        ID          uint      `gorm:"primary_key" json:"id"`
        Name        string    `json:"name"`
        Address 	string    `json:"address"`
		UserID 		uint 		`json:"userID"`
        User      	User  	 `json:"-"`
        Transaction      	[]Transaction  	 `json:"-"`
    }
)

func ExtractCustomer(id uint, db *gorm.DB) (uint, error) {
    // get db from gin context
    var customer Customer
    if err := db.Where("user_id = ?", id).First(&customer).Error; err != nil {
        return 0, err
    }

    idCustomer := customer.ID

    return idCustomer, nil
}