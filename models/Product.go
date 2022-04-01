package models


type (
    Product struct {
        ID          uint      `gorm:"primary_key" json:"id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
		Price 		int     `json:"price"`
		CategoryID  uint `json:"categoryID"`
		SellerID    uint `json:"sellerID"`
        Category      Category   `json:"-"`
        Seller      Seller   `json:"-"`
        Reviews      []Reviews   `json:"-"`
    }
)