package controller

import (
    "Final-Project/models"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type LoginInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// LoginUser godoc
// @Summary Login as as user.
// @Description Logging in to get jwt token to access admin or user api by roles.
// @Tags Auth
// @Param Body body LoginInput true "the body to login a user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input LoginInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    u := models.User{}

    u.Username = input.Username
    u.Password = input.Password

    token, err := models.LoginCheck(u.Username, u.Password, db)

    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
        return
    }

    user := map[string]string{
        "username": u.Username,
        "password":    u.Password,
    }

    c.JSON(http.StatusOK, gin.H{"message": "login success", "user": user, "token": token})

}

// Register godoc
// @Summary Register a user.
// @Description registering a user from public access.
// @Tags Auth
// @Param Body body RegisterInput true "the body to register a user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input RegisterInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    u := models.User{}

    u.Username = input.Username
    u.Password = input.Password
	u.Role = "Guest"

    _, err := u.SaveUser(db)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := map[string]string{
        "username": input.Username,
		"Role": u.Role,
    }

    c.JSON(http.StatusOK, gin.H{"message": "registration success", "user": user})

}