package router

import (
	"situkang/middleware"
	"situkang/models/entity"
	"situkang/provider"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func ConnectionRouter(router *gin.Engine, controller provider.ControllerProvider) {
	connectionController := controller.ProvideConnectionController()
	authController := controller.ProvideAuthController()
	userController := controller.ProvideUserController()
	homeController := controller.ProvideHomeController()
	categoryController := controller.ProvideCategoryController()
	workerPublicController := controller.ProvideWorkerPublicController()
	orderController := controller.ProvideOrderController()
	knowledgeController := controller.ProvideKnowledgeController()
	notificationController := controller.ProvideNotificationController()
	workerController := controller.ProvideWorkerController()

	// Swagger UI and OpenAPI documentation routes
	router.StaticFile("/docs/openapi.yaml", "./docs/openapi.yaml")
	router.StaticFile("/docs", "./docs/index.html")
	router.StaticFile("/swagger/openapi.yaml", "./docs/openapi.yaml")
	router.StaticFile("/swagger/index.html", "./docs/index.html")
	router.StaticFile("/swagger", "./docs/index.html")

	api := router.Group("/v1")
	api.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		api.POST("/connect", connectionController.Connect)
	}

	payments := api.Group("/payments")
	{
		payments.GET("/sandbox-checkout", orderController.SandboxCheckout)
		payments.POST("/sandbox-callback", orderController.SandboxCallback)
		payments.POST("/midtrans-webhook", orderController.HandleMidtransWebhook)
	}

	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.Refresh)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/reset-password", authController.ResetPassword)
		auth.POST("/logout", middleware.RequireAuth(), authController.Logout)
	}

	api.GET("/categories", categoryController.ListCategories)
	api.GET("/categories/:category_id/services", categoryController.ListCategoryServices)

	api.GET("/knowledge/articles", knowledgeController.ListArticles)
	api.GET("/knowledge/articles/:article_id", knowledgeController.GetArticle)
	api.GET("/knowledge/faq", knowledgeController.ListFAQ)

	authenticated := api.Group("/")
	authenticated.Use(middleware.RequireAuth())
	{
		authenticated.GET("/users/me", userController.GetMe)
		authenticated.PUT("/users/me", userController.UpdateMe)
		authenticated.PUT("/users/me/avatar", userController.UpdateAvatar)
		authenticated.PUT("/users/me/location", userController.UpdateLocation)

		authenticated.GET("/notifications", notificationController.ListNotifications)
		authenticated.PATCH("/notifications/:notification_id/read", notificationController.MarkRead)
		authenticated.PATCH("/notifications/read-all", notificationController.MarkAllRead)
	}

	userOnly := api.Group("/")
	userOnly.Use(middleware.RequireAuth(), middleware.RequireRoles(entity.UserRoleUser))
	{
		userOnly.GET("/home", homeController.GetUserHome)
		userOnly.GET("/workers/nearby", workerPublicController.ListNearby)
		userOnly.GET("/workers/search", workerPublicController.Search)
		userOnly.GET("/workers/:worker_id", workerPublicController.GetDetail)
		userOnly.GET("/workers/:worker_id/reviews", workerPublicController.GetReviews)
		userOnly.GET("/workers/:worker_id/services", workerPublicController.GetServices)

		userOnly.POST("/orders", orderController.CreateOrder)
		userOnly.GET("/orders", orderController.ListOrders)
		userOnly.GET("/orders/:order_id", orderController.GetOrderDetail)
		userOnly.POST("/orders/:order_id/cancel", orderController.CancelOrder)
		userOnly.GET("/orders/:order_id/tracking", orderController.GetTracking)
		userOnly.GET("/orders/:order_id/tracking/location", orderController.GetTrackingLocation)
		userOnly.GET("/orders/:order_id/purchases", orderController.ListPurchases)
		userOnly.GET("/orders/:order_id/purchases/:purchase_id", orderController.GetPurchaseDetail)
		userOnly.PATCH("/orders/:order_id/purchases/:purchase_id/approve", orderController.ApprovePurchase)
		userOnly.PATCH("/orders/:order_id/purchases/:purchase_id/reject", orderController.RejectPurchase)
		userOnly.PATCH("/orders/:order_id/purchases/:purchase_id/clarify", orderController.ClarifyPurchase)
		userOnly.PATCH("/orders/:order_id/purchases/bulk-approve", orderController.BulkApprovePurchases)
		userOnly.GET("/orders/:order_id/chat/messages", orderController.ListChatMessages)
		userOnly.POST("/orders/:order_id/chat/messages", orderController.SendChatMessage)
		userOnly.PATCH("/orders/:order_id/chat/read", orderController.MarkChatRead)
		userOnly.GET("/chats", orderController.ListChats)
		userOnly.POST("/orders/:order_id/rating", orderController.CreateRating)
		userOnly.GET("/orders/:order_id/rating", orderController.GetRating)
		userOnly.GET("/orders/:order_id/invoice", orderController.GetInvoice)
		userOnly.POST("/orders/:order_id/payment", orderController.CreatePayment)
		userOnly.POST("/orders/:order_id/payment/sync", orderController.SyncMidtransPayment)
		userOnly.GET("/orders/:order_id/invoice/pdf", orderController.DownloadInvoicePDF)
	}

	workerOnly := api.Group("/worker")
	workerOnly.Use(middleware.RequireAuth(), middleware.RequireRoles(entity.UserRoleWorker))
	{
		workerOnly.GET("/profile", workerController.GetProfile)
		workerOnly.PUT("/profile", workerController.UpdateProfile)
		workerOnly.PUT("/profile/cover-photo", workerController.UpdateCoverPhoto)
		workerOnly.POST("/profile/verification", workerController.SubmitVerification)
		workerOnly.GET("/profile/verification", workerController.GetVerification)

		workerOnly.GET("/home", workerController.GetHome)
		workerOnly.PATCH("/availability", workerController.UpdateAvailability)

		workerOnly.GET("/orders/incoming", workerController.ListIncomingOrders)
		workerOnly.GET("/orders/incoming/:order_id", workerController.GetIncomingOrderDetail)
		workerOnly.POST("/orders/:order_id/accept", workerController.AcceptOrder)
		workerOnly.POST("/orders/:order_id/reject", workerController.RejectOrder)
		workerOnly.GET("/orders", workerController.ListOrders)
		workerOnly.GET("/orders/:order_id", workerController.GetOrderDetail)
		workerOnly.PATCH("/orders/:order_id/status", workerController.UpdateOrderStatus)
		workerOnly.POST("/orders/:order_id/generate-invoice", workerController.GenerateInvoice)

		workerOnly.POST("/orders/:order_id/purchases", workerController.AddPurchase)
		workerOnly.POST("/orders/:order_id/purchases/ai-process", workerController.AIProcessPurchase)
		workerOnly.POST("/orders/:order_id/purchases/receipt-scan", workerController.ReceiptScanPurchase)
		workerOnly.PUT("/orders/:order_id/purchases/:purchase_id", workerController.UpdatePurchase)
		workerOnly.DELETE("/orders/:order_id/purchases/:purchase_id", workerController.DeletePurchase)
		workerOnly.POST("/orders/:order_id/purchases/:purchase_id/submit", workerController.SubmitPurchase)
		workerOnly.POST("/orders/:order_id/purchases/bulk-submit", workerController.BulkSubmitPurchase)
		workerOnly.PATCH("/orders/:order_id/purchases/:purchase_id/clarify-response", workerController.ClarifyPurchaseResponse)

		workerOnly.GET("/orders/:order_id/chat/messages", workerController.ListChatMessages)
		workerOnly.POST("/orders/:order_id/chat/messages", workerController.SendChatMessage)
		workerOnly.PATCH("/orders/:order_id/chat/read", workerController.MarkChatRead)
		workerOnly.GET("/chats", workerController.ListChats)

		workerOnly.POST("/orders/:order_id/customer-rating", workerController.CreateCustomerRating)
		workerOnly.GET("/orders/:order_id/customer-rating", workerController.GetCustomerRating)

		workerOnly.GET("/history", workerController.GetHistory)
		workerOnly.GET("/statistics", workerController.GetStatistics)

		workerOnly.GET("/wallet", workerController.GetWallet)
		workerOnly.GET("/wallet/transactions", workerController.ListWalletTransactions)
		workerOnly.POST("/wallet/withdraw", workerController.Withdraw)

		workerOnly.PUT("/location", workerController.UpdateLocation)
	}
}
