package controller

import (
    "net/http"

    "Final-Project/models"
    "Final-Project/utils/token"
    // "strconv"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type ProductInput struct {
    Name        string    `json:"name"`
	Description string    `json:"description"`
	Price 		int     `json:"price"`
	CategoryID  uint `json:"categoryID"`
}

// GetAllProduct godoc
// @Summary Get all Product.
// @Description Get a list of Product.
// @Tags Product
// @Produce json
// @Success 200 {object} []models.Product
// @Router /product [get]
func GetAllProduct(c *gin.Context) {
    // get db from gin context
    db := c.MustGet("db").(*gorm.DB)
    var product []models.Product
    db.Find(&product)

    c.JSON(http.StatusOK, gin.H{"data": product})
}

// GetProductById godoc
// @Summary Get Product By Id.
// @Description Get an Product by id.
// @Tags Product
// @Produce json
// @Param id path string true "Product id"
// @Success 200 {object} models.Product
// @Router /product/{id} [get]
func GetProductById(c *gin.Context) { // Get model if exist
    var product models.Product

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": product})
}

// CreateProduct godoc
// @Summary Create a new Product.
// @Description Creating a new Product.
// @Tags Product
// @Param Body body ProductInput true "the body to create a new Product"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Product
// @Router /product [post]
func CreateProduct(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input ProductInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)

    role, _ := models.ExtractRole(loggedId, db)
    idSeller, _ := models.ExtractSeller(loggedId, db)

    if role != "Seller" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "only seller could create a new product"})
        return
    }

    product := models.Product{Name: input.Name, 
        Description: input.Description, 
        Price: input.Price, 
        CategoryID: input.CategoryID, 
        SellerID: idSeller}
    db.Create(&product)

    c.JSON(http.StatusOK, gin.H{"data": product})

}


// UpdateProduct godoc
// @Summary Update Product.
// @Description Update Product by id.
// @Tags Product
// @Produce json
// @Param id path string true "Product id"
// @Param Body body CustomerInput true "the body to update product"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.Product
// @Router /product/{id} [patch]
func UpdateProduct(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var product models.Product
    if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)
	idSeller, _ := models.ExtractSeller(loggedId, db)

    // idString := strconv.FormatUint(uint64(loggedId), 10)

	if idSeller != product.SellerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't edit this data!"})
        return
	}

    // Validate input
    var input ProductInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedInput models.Product
    updatedInput.Name = input.Name
    updatedInput.Description = input.Description
    updatedInput.Price = input.Price
    updatedInput.CategoryID = input.CategoryID

    db.Model(&product).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteProduct godoc
// @Summary Delete one Product.
// @Description Delete a Product by id.
// @Tags Product
// @Produce json
// @Param id path string true "Product id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /product/{id} [delete]
func DeleteProduct(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var product models.Product

	loggedId, _ := token.ExtractTokenID(c)
	idSeller, _ := models.ExtractSeller(loggedId, db)

	
    if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

	if idSeller != product.SellerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't delete this data!"})
		return
	}

    db.Delete(&product)

    c.JSON(http.StatusOK, gin.H{"data": true})
}