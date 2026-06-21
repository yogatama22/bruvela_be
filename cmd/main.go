package main

import (
	"bruvela-backend/config"
	"bruvela-backend/internal/handler"
	"bruvela-backend/internal/middleware"
	"bruvela-backend/internal/repository"
	"bruvela-backend/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	ingredientRepo := repository.NewIngredientRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	batchRepo := repository.NewBatchRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	shippingTypeRepo := repository.NewShippingTypeRepository(db)
	paymentStatusRepo := repository.NewPaymentStatusRepository(db)
	orderStatusRepo := repository.NewOrderStatusRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	stockLogRepo := repository.NewStockLogRepository(db)
	companySettingRepo := repository.NewCompanySettingRepository(db)
	reportHandler := handler.NewReportHandler(orderRepo, financeRepo, ingredientRepo, recipeRepo, productRepo)

	authHandler := handler.NewAuthHandler(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	productHandler := handler.NewProductHandler(productRepo, recipeRepo)
	ingredientHandler := handler.NewIngredientHandler(ingredientRepo, recipeRepo, orderRepo)
	recipeHandler := handler.NewRecipeHandler(recipeRepo, productRepo, ingredientRepo)
	orderHandler := handler.NewOrderHandler(orderRepo, recipeRepo, ingredientRepo, stockLogRepo, db)
	financeHandler := handler.NewFinanceHandler(financeRepo)
	dashboardHandler := handler.NewDashboardHandler(orderRepo, ingredientRepo, batchRepo)
	shippingTypeHandler := handler.NewShippingTypeHandler(shippingTypeRepo)
	paymentStatusHandler := handler.NewPaymentStatusHandler(paymentStatusRepo)
	orderStatusHandler := handler.NewOrderStatusHandler(orderStatusRepo)
	batchHandler := handler.NewBatchHandler(batchRepo, orderRepo, financeRepo, recipeRepo)
	purchaseHandler := handler.NewPurchaseHandler(purchaseRepo, ingredientRepo, financeRepo, stockLogRepo, db)
	stockLogHandler := handler.NewStockLogHandler(stockLogRepo)
	companySettingHandler := handler.NewCompanySettingHandler(companySettingRepo)

	gin.SetMode(cfg.Server.GinMode)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(cfg.JWT.Secret), authHandler.GetMe)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			products := protected.Group("/products")
			{
				products.GET("", productHandler.GetAll)
				products.GET("/:id", productHandler.GetByID)
				products.POST("", productHandler.Create)
				products.PUT("/:id", productHandler.Update)
				products.DELETE("/:id", productHandler.Delete)
				products.GET("/:id/recipe", productHandler.GetRecipe)
			}

			ingredients := protected.Group("/ingredients")
			{
				ingredients.GET("", ingredientHandler.GetAll)
				ingredients.GET("/estimation", ingredientHandler.GetWithEstimation)
				ingredients.GET("/alerts", ingredientHandler.GetLowStock)
				ingredients.GET("/:id", ingredientHandler.GetByID)
				ingredients.POST("", ingredientHandler.Create)
				ingredients.PUT("/:id", ingredientHandler.Update)
				ingredients.DELETE("/:id", ingredientHandler.Delete)

				purchases := protected.Group("/ingredient-purchases")
				{
					purchases.GET("", purchaseHandler.GetAll)
					purchases.GET("/:id", purchaseHandler.GetByID)
					purchases.POST("", purchaseHandler.Create)
					purchases.PUT("/:id", purchaseHandler.Update)
					purchases.DELETE("/:id", purchaseHandler.Delete)
				}
			}

			recipes := protected.Group("/recipes")
			{
				recipes.GET("", recipeHandler.GetAllProductsWithHPP)
				recipes.GET("/product/:product_id", recipeHandler.GetByProductID)
				recipes.POST("", recipeHandler.Create)
				recipes.POST("/calculator", recipeHandler.CalculateProduction)
				recipes.PUT("/:id", recipeHandler.Update)
				recipes.DELETE("/:id", recipeHandler.Delete)
				recipes.DELETE("/product/:product_id", recipeHandler.DeleteByProductID)
			}

			orders := protected.Group("/orders")
			{
				orders.GET("", orderHandler.GetAll)
				orders.GET("/:id", orderHandler.GetByID)
				orders.POST("", orderHandler.Create)
				orders.PUT("/:id", orderHandler.Update)
				orders.PATCH("/:id/status", orderHandler.UpdateStatus)
				orders.PATCH("/:id/pay", orderHandler.UpdatePayStatus)
				orders.DELETE("/:id", orderHandler.Delete)
			}

			finance := protected.Group("/finance")
			{
				finance.GET("/journal", financeHandler.GetJournalEntries)
				finance.POST("/journal", financeHandler.CreateJournalEntry)
				finance.GET("/summary", financeHandler.GetSummary)
			}

			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/summary", dashboardHandler.GetSummary)
			}

			shippingTypes := protected.Group("/shipping-types")
			{
				shippingTypes.GET("", shippingTypeHandler.GetAll)
				shippingTypes.GET("/:id", shippingTypeHandler.GetByID)
				shippingTypes.POST("", shippingTypeHandler.Create)
				shippingTypes.PUT("/:id", shippingTypeHandler.Update)
				shippingTypes.DELETE("/:id", shippingTypeHandler.Delete)
			}

			paymentStatuses := protected.Group("/payment-statuses")
			{
				paymentStatuses.GET("", paymentStatusHandler.GetAll)
				paymentStatuses.GET("/:id", paymentStatusHandler.GetByID)
				paymentStatuses.POST("", paymentStatusHandler.Create)
				paymentStatuses.PUT("/:id", paymentStatusHandler.Update)
				paymentStatuses.DELETE("/:id", paymentStatusHandler.Delete)
			}

			orderStatuses := protected.Group("/order-statuses")
			{
				orderStatuses.GET("", orderStatusHandler.GetAll)
				orderStatuses.GET("/:id", orderStatusHandler.GetByID)
				orderStatuses.POST("", orderStatusHandler.Create)
				orderStatuses.PUT("/:id", orderStatusHandler.Update)
				orderStatuses.DELETE("/:id", orderStatusHandler.Delete)
			}

			batches := protected.Group("/batches")
			{
				batches.GET("", batchHandler.GetAll)
				batches.GET("/active", batchHandler.GetActive)
				batches.GET("/:id", batchHandler.GetByID)
				batches.GET("/:id/summary", batchHandler.GetSummary)
				batches.POST("", batchHandler.Create)
				batches.PUT("/:id", batchHandler.Update)
				batches.PATCH("/:id/activate", batchHandler.SetActive)
				batches.PATCH("/:id/close", batchHandler.Close)
				batches.DELETE("/:id", batchHandler.Delete)
			}

			stockLogs := protected.Group("/stock-logs")
			{
				stockLogs.GET("", stockLogHandler.GetAll)
			}

			reports := protected.Group("/reports")
			{
				reports.GET("/orders", reportHandler.ExportOrders)
				reports.GET("/finance", reportHandler.ExportFinance)
				reports.GET("/inventory", reportHandler.ExportInventory)
				reports.GET("/hpp", reportHandler.ExportHPP)
			}

			companySettings := protected.Group("/company-settings")
			{
				companySettings.GET("", companySettingHandler.Get)
				companySettings.PUT("", companySettingHandler.Upsert)
			}
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
