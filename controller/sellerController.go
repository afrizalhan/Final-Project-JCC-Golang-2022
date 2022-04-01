package controller

import (
    "net/http"

    "Final-Project/models"
    "Final-Project/utils/token"

    // "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type SellerInput struct {
    Name        string `json:"name"`
    Address string `json:"address"`
}

// GetAllSeller godoc
// @Summary Get all Seller.
// @Description Get a list of Seller
// @Tags Seller
// @Produce json
// @Success 200 {object} []models.Seller
// @Router /seller [get]
func GetAllSeller(c *gin.Context) {
    // get db from gin context
    db := c.MustGet("db").(*gorm.DB)
    var seller []models.Seller
    db.Find(&seller)

    c.JSON(http.StatusOK, gin.H{"data": seller})
}

// GetSellerById godoc
// @Summary Get Seller by Id.
// @Description Get an Seller by id.
// @Tags Seller
// @Produce json
// @Param id path string true "Seller id"
// @Success 200 {object} models.Seller
// @Router /seller/{id} [get]
func GetSellerById(c *gin.Context) { // Get model if exist
    var seller models.Seller

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&seller).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": seller})
}

// RegisterAsSeller godoc
// @Summary Register user as a Seller.
// @Description Creating a new Seller.
// @Tags Seller
// @Param Body body SellerInput true "the body to create a new Seller"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Seller
// @Router /seller [post]
func RegisterAsSeller(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input SellerInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)

    role, _ := models.ExtractRole(loggedId, db)

    if role != "Guest" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "you already registered as another role"})
        c.Abort()
        return
    }

    seller := models.Seller{Name: input.Name, Address: input.Address, UserID: loggedId}
    db.Create(&seller)

	var user models.User
    if err := db.Where("id = ?", loggedId).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

	var updatedInput models.User
    updatedInput.Role = "Seller"

    db.Model(&user).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": seller})

}

// UpdateSeller godoc
// @Summary Update Seller.
// @Description Update seller by id.
// @Tags Seller
// @Produce json
// @Param id path string true "Seller id"
// @Param Body body CustomerInput true "the body to update seller"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.Seller
// @Router /seller/{id} [patch]
func UpdateSeller(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var seller models.Seller
    if err := db.Where("id = ?", c.Param("id")).First(&seller).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }


	loggedId, _ := token.ExtractTokenID(c)

    // idString := strconv.FormatUint(uint64(loggedId), 10)

	if loggedId != seller.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't edit this data!"})
        return
	}

    // Validate input
    var input SellerInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedInput models.Seller
    updatedInput.Name = input.Name
    updatedInput.Address = input.Address

    db.Model(&seller).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": seller})
}

// DeleteSeller godoc
// @Summary Delete one Seller.
// @Description Delete a Seller by id.
// @Tags Seller
// @Produce json
// @Param id path string true "Seller id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /seller/{id} [delete]
func DeleteSeller(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var seller models.Seller

	loggedId, _ := token.ExtractTokenID(c)

    // idString := strconv.FormatUint(uint64(loggedId), 10)

    
    if err := db.Where("id = ?", c.Param("id")).First(&seller).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }
    if loggedId != seller.UserID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You can't delete this data!"})
        return
    }

    db.Delete(&seller)

    c.JSON(http.StatusOK, gin.H{"data": true})
}