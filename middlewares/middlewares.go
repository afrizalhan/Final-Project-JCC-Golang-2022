package middlewares

import (
    "net/http"

    "Final-Project/utils/token"
    "Final-Project/models"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        err := token.TokenValid(c)
        if err != nil {
            c.String(http.StatusUnauthorized, err.Error())
            c.Abort()
            return
        }
        c.Next()
    }
}

func JwtAuthAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        err := token.TokenValid(c)
        if err != nil {
            c.String(http.StatusUnauthorized, err.Error())
            c.Abort()
            return
        }
        db := c.MustGet("db").(*gorm.DB)
        loggedId, _ := token.ExtractTokenID(c)

        role, _ := models.ExtractRole(loggedId, db)

        if role != "Admin" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "this operation only open for user with Admin role"})
            c.Abort()
            return
        }
        c.Next()
    }
}

