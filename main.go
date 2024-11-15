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
	initializers.InitRedis()
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	v1 := r.Group("/api/v1")
	{
		productRoutes := v1.Group("/products")
		{
			productRoutes.GET("/", controllers.GetProducts)                                                                // GET all products
			productRoutes.POST("/", middleware.Authorization(config.CreateProduct), controllers.InsertProduct)             // CREATE a new product
			productRoutes.GET("/id/:id", controllers.GetProductById)                                                       // GET product by id
			productRoutes.PUT("/id/:id", middleware.Authorization(config.EditProduct), controllers.EditProductById)        // EDIT the product
			productRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteProduct), controllers.DeleteProductById) // DELETE the product
			productRoutes.GET("/search/:phrase", controllers.GetProductsBySearch)                                          // SEARCH for product
		}

		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", controllers.Login)       // LOGIN
			authRoutes.POST("/register", controllers.Register) // REGISTER
			authRoutes.POST("/logout", controllers.Logout)     // LOGOUT
		}

		variantRoutes := v1.Group("/variants")
		{
			variantRoutes.GET("/", middleware.Authorization(config.ReadVariant))
			variantRoutes.POST("/", middleware.Authorization(config.CreateVariant))
			variantRoutes.GET("/id/:id", middleware.Authorization(config.ReadVariant))
			variantRoutes.PUT("/id/:id", middleware.Authorization(config.EditVariant))
			variantRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteVariant))
		}

		categoryRoutes := v1.Group("/categories")
		{
			categoryRoutes.GET("/")
			categoryRoutes.POST("/", middleware.Authorization(config.CreateCategory))
			categoryRoutes.GET("/id/:id")
			categoryRoutes.PUT("/id/:id", middleware.Authorization(config.EditCategory))
			categoryRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteCategory))
		}

		subcategoryRoutes := v1.Group("/subcategories")
		{
			subcategoryRoutes.GET("")
			subcategoryRoutes.POST("/", middleware.Authorization(config.CreateSubcategory))
			subcategoryRoutes.GET("/id/:id")
			subcategoryRoutes.PUT("/id/:id", middleware.Authorization(config.EditSubcategory))
			subcategoryRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteSubcategory))
		}
	}

	r.Run()
}
