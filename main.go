package main

import (
	"litstore/api/config"
	"litstore/api/controllers"
	"litstore/api/initializers"
	"litstore/api/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	v1 := r.Group("/api/v1")
	{
		productRoutes := v1.Group("/products")
		{
			productRoutes.GET("/", controllers.GetProducts)
			productRoutes.POST("/", middleware.Authorization(config.CreateProduct), controllers.InsertProduct)

			productRoutes.GET("/id/:id", controllers.GetProductById)
			productRoutes.PUT("/id/:id", controllers.EditProductById)
			productRoutes.DELETE("/id/:id", controllers.DeleteProductById)

			productRoutes.GET("/search/:phrase", controllers.GetProductsBySearch)
		}

		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", controllers.Login)
			authRoutes.POST("/register", controllers.Register)
		}

		variantRoutes := v1.Group("/variants")
		{
			variantRoutes.GET("/")
			variantRoutes.POST("/")
		}
	}

	r.Run()
}
