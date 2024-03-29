package models

import (
    "html"
    "strings"
    "time"

    "Final-Project/utils/token"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
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

func VerifyPassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string, db *gorm.DB) (string, error) {

    var err error

    u := User{}

    err = db.Model(User{}).Where("username = ?", username).Take(&u).Error

    if err != nil {
        return "", err
    }

    err = VerifyPassword(password, u.Password)

    if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
        return "", err
    }

    token, err := token.GenerateToken(u.ID)

    if err != nil {
        return "", err
    }

    return token, nil

}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
    //turn password into hash
    hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if errPassword != nil {
        return &User{}, errPassword
    }
    u.Password = string(hashedPassword)
    //remove spaces in username
    u.Username = html.EscapeString(strings.TrimSpace(u.Username))

    var err error = db.Create(&u).Error
    if err != nil {
        return &User{}, err
    }
    return u, nil
}


func ExtractRole(id uint, db *gorm.DB) (string, error) {
    // get db from gin context
    var user User
    if err := db.Where("id = ?", id).First(&user).Error; err != nil {
        return "not found", err
    }

    role := user.Role

    return role, nil
}

