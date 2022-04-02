package controller

import (
    "net/http"

    "Final-Project/models"
    "Final-Project/utils/token"
    // "strconv"
    "time"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type TransactionInput struct {
    ProductID 	uint 	`json:"productID"`
	Description string    `json:"description"`
	Quantity 	int 			`json:"quantity"`
	Status   	string		 `json:"status"`
}

// GetAllTransaction godoc
// @Summary Get all Transaction (Admin Only).
// @Description Get a list of Transaction.
// @Tags Transaction
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.Transaction
// @Router /transaction [get]
func GetAllTransaction(c *gin.Context) {
    // get db from gin context
    db := c.MustGet("db").(*gorm.DB)
    loggedId, _ := token.ExtractTokenID(c)

    role, _ := models.ExtractRole(loggedId, db)

    if role != "Admin" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "only admin could access this data"})
        return
    }
    var transaction []models.Transaction
    db.Find(&transaction)
    

    c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// GetTransactionById godoc
// @Summary Get Transaction By Id (Admin and Customer in Transaction Only).
// @Description Get an Transaction by id.
// @Tags Transaction
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Transaction id"
// @Success 200 {object} models.Transaction
// @Router /transaction/{id} [get]
func GetTransactionById(c *gin.Context) { // Get model if exist
    var transaction models.Transaction
    var product models.Product
    db := c.MustGet("db").(*gorm.DB)

    loggedId, _ := token.ExtractTokenID(c)

    role, _ := models.ExtractRole(loggedId, db)
    var idCustomer uint
    if role == "Customer"{
        idCustomer, _ = models.ExtractCustomer(loggedId, db)
    }
    idSeller, _ := models.ExtractSeller(loggedId, db)
    
    if err := db.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }
    if err := db.Where("id = ?", transaction.ProductID).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
        return
    }
    if role == "Admin" || idCustomer == transaction.CustomerID || idSeller == product.SellerID{
        c.JSON(http.StatusOK, gin.H{"data": transaction})
        return
    }
    c.JSON(http.StatusBadRequest, gin.H{"error": "only admin and related users could access this data"})
}

// CreateTransaction godoc
// @Summary Create a new Transaction (Customer Only).
// @Description Creating a new Transaction.
// @Tags Transaction
// @Param Body body TransactionInput true "the body to create a new Product"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Transaction
// @Router /transaction [post]
func CreateTransaction(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input TransactionInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)

    role, _ := models.ExtractRole(loggedId, db)

    if role != "Customer" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "only customer could create a new transaction"})
        return
    }
	
	var product models.Product
    if err := db.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
        return
    }

	customerId, _ := models.ExtractCustomer(loggedId, db)

    transaction := models.Transaction{ProductID: input.ProductID, 
		Description: input.Description, 
		Quantity: input.Quantity, 
		TotalPrice: input.Quantity * product.Price, 
		Status: input.Status, 
		CustomerID: customerId, 
		SellerID: product.SellerID}
    db.Create(&transaction)

    c.JSON(http.StatusOK, gin.H{"data": transaction})

}


// UpdateTransaction godoc
// @Summary Update Transaction (Customer and Seller in Transaction Only).
// @Description Update Transaction by id.
// @Tags Transaction
// @Produce json
// @Param id path string true "Transaction id"
// @Param Body body TransactionInput true "the body to update transaction"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.Transaction
// @Router /transaction/{id} [patch]
func UpdateTransaction(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var transaction models.Transaction
    if err := db.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    var product models.Product
    if err := db.Where("id = ?", transaction.ProductID).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }
    loggedId, _ := token.ExtractTokenID(c)

	role, _ := models.ExtractRole(loggedId, db)
    var idCustomer uint
    if role == "Customer"{
        idCustomer, _ = models.ExtractCustomer(loggedId, db)
    }
    idSeller, _ := models.ExtractSeller(loggedId, db)

	if idCustomer != transaction.CustomerID && product.SellerID != idSeller {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't edit this data!"})
        return
	}

    // Validate input
    var input TransactionInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedInput models.Transaction
	productPrice := transaction.TotalPrice / transaction.Quantity
    updatedInput.Description = input.Description
    updatedInput.Quantity = input.Quantity
	updatedInput.TotalPrice = input.Quantity * productPrice
    updatedInput.Status = input.Status
    updatedInput.UpdatedAt = time.Now()

    db.Model(&transaction).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// DeleteTransaction godoc
// @Summary Delete one Transaction (Admin Only).
// @Description Delete a Transaction by id.
// @Tags Transaction
// @Produce json
// @Param id path string true "Transaction id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /transaction/{id} [delete]
func DeleteTransaction(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var transaction models.Product

	loggedId, _ := token.ExtractTokenID(c)
    role, _ := models.ExtractRole(loggedId, db)

	if role != "Admin" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "only admin could delete transaction data!"})
        return
    }
    if err := db.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    db.Delete(&transaction)

    c.JSON(http.StatusOK, gin.H{"data": true})
}