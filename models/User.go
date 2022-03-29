package models

import (
    "time"
)

type (
    User struct {
        ID          uint      `gorm:"primary_key" json:"id"`
        Username    string    `json:"username"`
        Password 	string    `json:"password"`
		Role 		string 	  `json:"role"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
        Customer      []Customer   `json:"-"`
        Seller      []Seller   `json:"-"`
		Reviews		 []Reviews   `json:"-"`
    }
)