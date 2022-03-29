package models

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