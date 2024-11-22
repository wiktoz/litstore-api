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
			productRoutes.GET("/slug/:slug", controllers.GetProductBySlug)                                                 // GET product by slug
			productRoutes.PUT("/id/:id", middleware.Authorization(config.EditProduct), controllers.EditProductById)        // EDIT the product
			productRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteProduct), controllers.DeleteProductById) // DELETE the product
			productRoutes.GET("/search/:phrase", controllers.GetProductsBySearch)                                          // SEARCH for product
		}

		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", controllers.Login)       // LOGIN
			authRoutes.POST("/register", controllers.Register) // REGISTER
			authRoutes.POST("/logout", controllers.Logout)     // LOGOUT
			authRoutes.POST("/password/forgot")
			authRoutes.POST("/password/change/:token")
		}

		variantRoutes := v1.Group("/variants")
		{
			variantRoutes.GET("/", middleware.Authorization(config.ReadVariant), controllers.GetVariants)
			variantRoutes.POST("/", middleware.Authorization(config.CreateVariant), controllers.InsertVariant)
			variantRoutes.GET("/id/:id", middleware.Authorization(config.ReadVariant), controllers.GetVariantById)
			variantRoutes.PUT("/id/:id", middleware.Authorization(config.EditVariant), controllers.EditVariantById)
			variantRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteVariant), controllers.DeleteVariantById)
		}

		categoryRoutes := v1.Group("/categories")
		{
			categoryRoutes.GET("/", controllers.GetCategories)
			categoryRoutes.POST("/", middleware.Authorization(config.CreateCategory), controllers.InsertCategory)
			categoryRoutes.GET("/id/:id", controllers.GetCategoryById)
			categoryRoutes.GET("/slug/:slug", controllers.GetCategoryBySlug)
			categoryRoutes.PUT("/id/:id", middleware.Authorization(config.EditCategory), controllers.EditCategoryById)
			categoryRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteCategory), controllers.DeleteCategoryById)
		}

		subcategoryRoutes := v1.Group("/subcategories")
		{
			subcategoryRoutes.GET("/", controllers.GetSubcategories)
			subcategoryRoutes.POST("/", middleware.Authorization(config.CreateSubcategory), controllers.InsertSubcategory)
			subcategoryRoutes.GET("/id/:id", controllers.GetSubcategoryById)
			subcategoryRoutes.GET("/slug/:slug", controllers.GetSubcategoryBySlug)
			subcategoryRoutes.PUT("/id/:id", middleware.Authorization(config.EditSubcategory), controllers.EditSubcategoryById)
			subcategoryRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteSubcategory), controllers.DeleteSubcategoryById)
		}

		userRoutes := v1.Group("/users")
		{
			userRoutes.GET("/me", middleware.Authorization(""), controllers.GetUserSelf)
			userRoutes.GET("/", middleware.Authorization(config.ReadUser), controllers.GetUsers)
			userRoutes.GET("/id/:id", middleware.Authorization(config.ReadUser), controllers.GetUserById)
			userRoutes.GET("/search/:phrase", middleware.Authorization(config.ReadUser), controllers.GetUsersBySearch)
			userRoutes.PUT("/id/:id", middleware.Authorization(config.EditUser), controllers.EditUserById)
			userRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteUser), controllers.DeleteUserById)
		}
	}

	r.Run()
}
