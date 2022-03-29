package models

import (
    "time"
)

type (
    // AgeRatingCategory
    Transaction struct {
        ID          uint      `gorm:"primary_key" json:"id"`
		ProductID 	uint 	`json:"productID"`
        Description string    `json:"description"`
		Quantity 	int 			`json:"quantity"`
		TotalPrice	int			`json:"total_price"`
		Status   	string		 `json:"string"`
		SellerID 	uint 	`json:"sellerID"`
		CustomerID 	uint 	`json:"userID"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
        Product     Product   `json:"-"`
        Seller      Seller   `json:"-"`
        Customer    Customer   `json:"-"`
    }
)