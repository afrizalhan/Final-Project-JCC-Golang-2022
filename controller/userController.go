package controller

import (
    "net/http"

    "Final-Project/models"
    "Final-Project/utils/token"

    "time"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
    Username        string `json:"username"`
    Password string `json:"password"`
}

// GetAllUser godoc
// @Summary Get all User.
// @Description Get a list of User.
// @Tags User
// @Produce json
// @Success 200 {object} []models.User
// @Router /user [get]
func GetAllUser(c *gin.Context) {
    // get db from gin context
    db := c.MustGet("db").(*gorm.DB)
    var user []models.User
    db.Find(&user)

    c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetUserById godoc
// @Summary Get User by Id.
// @Description Get an User by id.
// @Tags User
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} models.User
// @Router /user/{id} [get]
func GetUserById(c *gin.Context) { // Get model if exist
    var user models.User

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser godoc
// @Summary Update User (Referred User Only).
// @Description Update User by id.
// @Tags User
// @Produce json
// @Param id path string true "User id"
// @Param Body body UserInput true "the body to update user"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.User
// @Router /user/{id} [patch]
func UpdateUser(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var user models.User
    if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }


	loggedId, _ := token.ExtractTokenID(c)

	if loggedId != user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't edit this data!"})
        return
	}

    // Validate input
    var input UserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedInput models.User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    updatedInput.Username = input.Username
    updatedInput.Password = string(hashedPassword)
    updatedInput.UpdatedAt = time.Now()

    db.Model(&user).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": user})

}

// DeleteUser godoc
// @Summary Delete one User (Admin and Referred User Only).
// @Description Delete a User by id.
// @Tags User
// @Produce json
// @Param id path string true "User id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /user/{id} [delete]
func DeleteUser(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var user models.User
    if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    loggedId, _ := token.ExtractTokenID(c)
    role, _ := models.ExtractRole(loggedId, db)

    
    if loggedId != user.ID && role != "Admin" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You can't delete this data!"})
        return
    }
	if user.Role == "Seller"{
        var seller models.Seller
		if err := db.Where("user_id = ?", c.Param("id")).First(&seller).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
		db.Delete(&seller)
	}
	if user.Role == "Customer"{
		var customer models.Customer
		if err := db.Where("user_id = ?", c.Param("id")).First(&customer).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
		db.Delete(&customer)
	}

    db.Delete(&user)

    c.JSON(http.StatusOK, gin.H{"data": true})
}