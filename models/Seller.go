package models


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