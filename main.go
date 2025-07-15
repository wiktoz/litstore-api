package main

import (
	"litstore/api/config"
	"litstore/api/controllers"
	"litstore/api/initializers"
	"litstore/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "litstore/api/docs"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.SyncDatabase()
	initializers.InitRedis()
	initializers.InitR2Client()
}

// @title       Litstore WebAPI
// @version     1.0
// @description E-commerce system API Docs
// @host        localhost:8000
// @BasePath    /api/v1

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Custom CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	v1 := r.Group("/api/v1")
	{
		productRoutes := v1.Group("/products")
		{
			productRoutes.GET("/all", controllers.GetProducts)                                                             // GET all products
			productRoutes.POST("/new", middleware.Authorization(config.CreateProduct), controllers.InsertProduct)          // CREATE a new product
			productRoutes.GET("/id/:id", controllers.GetProductById)                                                       // GET product by id
			productRoutes.GET("/slug/:slug", controllers.GetProductBySlug)                                                 // GET product by slug
			productRoutes.PUT("/id/:id", middleware.Authorization(config.EditProduct), controllers.EditProductById)        // EDIT the product
			productRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteProduct), controllers.DeleteProductById) // DELETE the product
			productRoutes.GET("/search/:phrase", controllers.GetProductsBySearch)                                          // SEARCH for product
			productRoutes.GET("/category/id/:category_id", controllers.GetProductsByCategory)                              // GET products by category
			productRoutes.GET("/subcategory/id/:subcategory_id", controllers.GetProductsBySubcategory)                     // GET products by subcategory
		}

		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", controllers.Login)                                 // LOGIN
			authRoutes.POST("/register", controllers.Register)                           // REGISTER
			authRoutes.POST("/logout", middleware.Authorization(""), controllers.Logout) // LOGOUT
			authRoutes.POST("/password/forgot", controllers.DemandResetPassword)         // FORGOT PASSWORD
			authRoutes.POST("/password/forgot/reset", controllers.ResetPassword)         // RESET PASSWORD
			authRoutes.POST("/email/verify", controllers.VerifyEmail)                    // VERIFY EMAIL
			authRoutes.POST("/email/resend", controllers.ResendVerificationEmail)        // RESEND VERIFICATION EMAIL
		}

		imagesRoutes := v1.Group("/images")
		{
			imagesRoutes.POST("/new", middleware.Authorization(config.CreateFile), controllers.UploadImages)
		}

		variantRoutes := v1.Group("/variants")
		{
			variantRoutes.GET("/all", middleware.Authorization(config.ReadVariant), controllers.GetVariants)
			variantRoutes.POST("/new", middleware.Authorization(config.CreateVariant), controllers.InsertVariant)
			variantRoutes.GET("/id/:id", middleware.Authorization(config.ReadVariant), controllers.GetVariantById)
			variantRoutes.PUT("/id/:id", middleware.Authorization(config.EditVariant), controllers.EditVariantById)
			variantRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteVariant), controllers.DeleteVariantById)
		}

		categoryRoutes := v1.Group("/categories")
		{
			categoryRoutes.GET("/all", controllers.GetCategories)
			categoryRoutes.POST("/new", middleware.Authorization(config.CreateCategory), controllers.InsertCategory)
			categoryRoutes.GET("/id/:id", controllers.GetCategoryById)
			categoryRoutes.GET("/slug/:slug", controllers.GetCategoryBySlug)
			categoryRoutes.PUT("/id/:id", middleware.Authorization(config.EditCategory), controllers.EditCategoryById)
			categoryRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteCategory), controllers.DeleteCategoryById)
		}

		subcategoryRoutes := v1.Group("/subcategories")
		{
			subcategoryRoutes.GET("/all", controllers.GetSubcategories)
			subcategoryRoutes.POST("/new", middleware.Authorization(config.CreateSubcategory), controllers.InsertSubcategory)
			subcategoryRoutes.GET("/id/:id", controllers.GetSubcategoryById)
			subcategoryRoutes.GET("/slug/:slug", controllers.GetSubcategoryBySlug)
			subcategoryRoutes.PUT("/id/:id", middleware.Authorization(config.EditSubcategory), controllers.EditSubcategoryById)
			subcategoryRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteSubcategory), controllers.DeleteSubcategoryById)
		}

		userRoutes := v1.Group("/users")
		{
			userRoutes.GET("/me", middleware.Authorization(""), controllers.GetUserSelf)
			userRoutes.GET("/all", middleware.Authorization(config.ReadUser), controllers.GetUsers)
			userRoutes.GET("/id/:id", middleware.Authorization(config.ReadUser), controllers.GetUserById)
			userRoutes.GET("/search/:phrase", middleware.Authorization(config.ReadUser), controllers.GetUsersBySearch)
			userRoutes.PUT("/id/:id", middleware.Authorization(config.EditUser), controllers.EditUserById)
			userRoutes.DELETE("/id/:id", middleware.Authorization(config.DeleteUser), controllers.DeleteUserById)
			userRoutes.POST("/password/change", middleware.Authorization(""), controllers.ChangePassword) // Change password

			userRoutes.POST("/address/new", middleware.Authorization(""), controllers.InsertUserAddress)   // Add new address
			userRoutes.GET("/address/all", middleware.Authorization(""), controllers.GetUserAddresses)     // Get all addresses
			userRoutes.DELETE("/address/:id", middleware.Authorization(""), controllers.DeleteUserAddress) // Delete address
		}
	}

	r.Run()
}
