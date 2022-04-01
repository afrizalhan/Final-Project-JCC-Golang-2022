package models

import (
    "time"
)

type (
    Reviews struct {
        ID          uint      `gorm:"primary_key" json:"id"`
        Rating        int    `json:"rating"`
		Comment string `json:"comment"`
		UserID 		uint `json:"userID"`
		ProductID 		uint `json:"productID"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
        User      User   `json:"-"`
        Product   Product   `json:"-"`
    }
)