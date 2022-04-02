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

type ReviewInput struct {
    Rating        	int    `json:"rating"`
	Comment 		string `json:"comment"`
	ProductID 		uint `json:"productID"`
}

// GetAllReview godoc
// @Summary Get all Review.
// @Description Get a list of Review.
// @Tags Review
// @Produce json
// @Success 200 {object} []models.Reviews
// @Router /review [get]
func GetAllReview(c *gin.Context) {
    // get db from gin context
    db := c.MustGet("db").(*gorm.DB)
    var review []models.Reviews
    db.Find(&review)

    c.JSON(http.StatusOK, gin.H{"data": review})
}

// GetReviewById godoc
// @Summary Get Review.
// @Description Get an Review by id.
// @Tags Review
// @Produce json
// @Param id path string true "Review id"
// @Success 200 {object} models.Reviews
// @Router /review/{id} [get]
func GetReviewById(c *gin.Context) { // Get model if exist
    var review models.Reviews

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&review).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": review})
}

// CreateReview godoc
// @Summary Create New Review (Customer Only).
// @Description Creating a new Review.
// @Tags Review
// @Param Body body ReviewInput true "the body to create a new Review"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Reviews
// @Router /review [post]
func CreateReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    // Validate input
    var input ReviewInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)

	role, _ := models.ExtractRole(loggedId, db)

    if role != "Customer" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "only customer could write review"})
        c.Abort()
        return
    }

    review := models.Reviews{Rating: input.Rating, Comment: input.Comment, UserID: loggedId, ProductID: input.ProductID}
    db.Create(&review)

    c.JSON(http.StatusOK, gin.H{"data": review})
}


// UpdateReview godoc
// @Summary Update Review (Customer Only).
// @Description Update Review by id.
// @Tags Review
// @Produce json
// @Param id path string true "Review id"
// @Param Body body ReviewInput true "the body to update review"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.Reviews
// @Router /review/{id} [patch]
func UpdateReview(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var review models.Reviews
    if err := db.Where("id = ?", c.Param("id")).First(&review).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

	loggedId, _ := token.ExtractTokenID(c)

    // idString := strconv.FormatUint(uint64(loggedId), 10)

	if loggedId != review.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't edit this data!"})
        return
	}

    // Validate input
    var input ReviewInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedInput models.Reviews
    updatedInput.Rating = input.Rating
    updatedInput.Comment = input.Comment
    updatedInput.UpdatedAt = time.Now()

    db.Model(&review).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"data": review})
}

// DeleteReview godoc
// @Summary Delete one Review (Customer that write the review Only).
// @Description Delete a Review by id.
// @Tags Review
// @Produce json
// @Param id path string true "Review id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /review/{id} [delete]
func DeleteReview(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var review models.Reviews

	loggedId, _ := token.ExtractTokenID(c)

    // idString := strconv.FormatUint(uint64(loggedId), 10)

    
    if err := db.Where("id = ?", c.Param("id")).First(&review).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }
    if loggedId != review.UserID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You can't delete this data!"})
        return
    }

    db.Delete(&review)

    c.JSON(http.StatusOK, gin.H{"data": true})
}