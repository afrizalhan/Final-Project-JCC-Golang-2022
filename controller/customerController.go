package controller

import (
    "net/http"

    "Final-Project/models"
    "Final-Project/utils/token"
    // "strconv"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type CustomerInput struct {
    Name        string `json:"name"`
    Address string `json:"address"`
}

// GetAllCustomer godoc
// @Summary Get all Customer.
// @Description Get a list of Customer.
// @Tags Customer
// @Produce json
// @Success 200 {object} []models.Customer
// @Router /customer [get]
func GetAllCustomer(c *gin.Context) {
    // get db from gin context
    db := c.MustGet("db").(*gorm.DB)
    var customer []models.Customer
    db.Find(&customer)

    c.JSON(http.StatusOK, gin.H{"data": customer})
}

// GetCustomerById godoc
// @Summary Get Customer.
// @Description Get an Customer by id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Success 200 {object} models.Customer
// @Router /customer/{id} [get]
func GetCustomerById(c *gin.Context) { // Get model if exist
    var customer models.Customer

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&customer).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": customer})
}

// RegisterAsCustomer godoc
// @Summary Register user as a Customer.
// @Description Creating a new Customer.
// @Tags Customer
// @Param Body body CustomerInput true "the body to create a new Customer"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Customer
// @Router /customer [post]
func RegisterAsCustomer(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input CustomerInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)

    role, _ := models.ExtractRole(loggedId, db)

    if role != "Guest" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "you already registered as another role"})
        return
    }

    customer := models.Customer{Name: input.Name, Address: input.Address, UserID: loggedId}
    db.Create(&customer)

    var user models.User
    if err := db.Where("id = ?", loggedId).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

	var updatedInput models.User
    updatedInput.Role = "Customer"

    db.Model(&user).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": customer})

}


// UpdateCustomer godoc
// @Summary Update Customer.
// @Description Update Customer by id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Param Body body CustomerInput true "the body to update customer"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.Customer
// @Router /customer/{id} [patch]
func UpdateCustomer(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var customer models.Customer
    if err := db.Where("id = ?", c.Param("id")).First(&customer).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)

    // idString := strconv.FormatUint(uint64(loggedId), 10)

	if loggedId != customer.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't edit this data!"})
        return
	}

    // Validate input
    var input CustomerInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedInput models.Customer
    updatedInput.Name = input.Name
    updatedInput.Address = input.Address

    db.Model(&customer).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": customer})
}

// DeleteCustomer godoc
// @Summary Delete one Customer.
// @Description Delete a Customer by id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /customer/{id} [delete]
func DeleteCustomer(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var customer models.Customer

	loggedId, _ := token.ExtractTokenID(c)

    // idString := strconv.FormatUint(uint64(loggedId), 10)

    
    if err := db.Where("id = ?", c.Param("id")).First(&customer).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }
    if loggedId != customer.UserID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You can't edit this data!"})
        return
    }

    db.Delete(&customer)

    c.JSON(http.StatusOK, gin.H{"data": true})
}