package routes

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "Final-Project/controller"
    "Final-Project/middlewares"

    swaggerFiles "github.com/swaggo/files"     // swagger embed files
    ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // set db to gin context
    r.Use(func(c *gin.Context) {
        c.Set("db", db)
    })

    r.POST("/register", controller.Register)
    r.POST("/login", controller.Login)

    r.GET("/seller", controller.GetAllSeller)
    r.GET("/customer", controller.GetAllCustomer)
    r.GET("/categories", controller.GetAllCategory)
    r.GET("/user", controller.GetAllUser)
    r.GET("/product", controller.GetAllProduct)
    r.GET("/review", controller.GetAllReview)
    r.GET("/seller/:id", controller.GetSellerById)
    r.GET("/customer/:id", controller.GetCustomerById)
    r.GET("/categories/:id", controller.GetCategoryById)
    r.GET("/user/:id", controller.GetUserById)
    r.GET("/product/:id", controller.GetProductById)
    r.GET("/review/:id", controller.GetReviewById)

    categoryMiddlewareRoute := r.Group("/categories")
    categoryMiddlewareRoute.Use(middlewares.JwtAuthAdmin())
    categoryMiddlewareRoute.POST("/", controller.CreateCategories)
    categoryMiddlewareRoute.PATCH("/:id", controller.UpdateCategory)
    categoryMiddlewareRoute.DELETE("/:id", controller.DeleteCategory)

    sellerMiddlewareRoute := r.Group("/seller")
    sellerMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
    sellerMiddlewareRoute.POST("/", controller.RegisterAsSeller)
    sellerMiddlewareRoute.PATCH("/:id", controller.UpdateSeller)
    sellerMiddlewareRoute.DELETE("/:id", controller.DeleteSeller)

	customerMiddlewareRoute := r.Group("/customer")
    customerMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
    customerMiddlewareRoute.POST("/", controller.RegisterAsCustomer)
    customerMiddlewareRoute.PATCH("/:id", controller.UpdateCustomer)
    customerMiddlewareRoute.DELETE("/:id", controller.DeleteCustomer)

	productMiddlewareRoute := r.Group("/product")
    productMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
    productMiddlewareRoute.POST("/", controller.CreateProduct)
    productMiddlewareRoute.PATCH("/:id", controller.UpdateProduct)
    productMiddlewareRoute.DELETE("/:id", controller.DeleteProduct)

	reviewMiddlewareRoute := r.Group("/review")
    reviewMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
    reviewMiddlewareRoute.POST("/", controller.CreateReview)
    reviewMiddlewareRoute.PATCH("/:id", controller.UpdateReview)
    reviewMiddlewareRoute.DELETE("/:id", controller.DeleteReview)

    transactionMiddlewareRoute := r.Group("/transaction")
    transactionMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
    transactionMiddlewareRoute.GET("/", controller.GetAllTransaction)
    transactionMiddlewareRoute.GET("/:id", controller.GetTransactionById)
    transactionMiddlewareRoute.POST("/", controller.CreateTransaction)
    transactionMiddlewareRoute.PATCH("/:id", controller.UpdateTransaction)
    transactionMiddlewareRoute.DELETE("/:id", controller.DeleteTransaction)

	userMiddlewareRoute := r.Group("/user")
    userMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
    userMiddlewareRoute.PATCH("/:id", controller.UpdateUser)
    userMiddlewareRoute.DELETE("/:id", controller.DeleteUser)

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return r
}