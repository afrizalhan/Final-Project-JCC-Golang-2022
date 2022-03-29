package models

type (
    // AgeRatingCategory
    Category struct {
        ID          uint      `gorm:"primary_key" json:"id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
        Product      []Product   `json:"-"`
    }
)