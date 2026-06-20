package provider

import "situkang/controllers"

type ControllerProvider interface {
	ProvideConnectionController() controllers.ConnectionController
	ProvideAuthController() controllers.AuthController
	ProvideUserController() controllers.UserController
	ProvideHomeController() controllers.HomeController
	ProvideCategoryController() controllers.CategoryController
	ProvideWorkerPublicController() controllers.WorkerPublicController
	ProvideOrderController() controllers.OrderController
	ProvideKnowledgeController() controllers.KnowledgeController
	ProvideNotificationController() controllers.NotificationController
	ProvideWorkerController() controllers.WorkerController
	ProvideFileController() controllers.FileController
}

type controllerProvider struct {
	connectionController   controllers.ConnectionController
	authController         controllers.AuthController
	userController         controllers.UserController
	homeController         controllers.HomeController
	categoryController     controllers.CategoryController
	workerPublicController controllers.WorkerPublicController
	orderController        controllers.OrderController
	knowledgeController    controllers.KnowledgeController
	notificationController controllers.NotificationController
	workerController       controllers.WorkerController
	fileController         controllers.FileController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	connectionController := controllers.NewConnectionController(servicesProvider.ProvideConnectionService())
	authController := controllers.NewAuthController(servicesProvider.ProvideAuthService())
	userController := controllers.NewUserController(servicesProvider.ProvideUserService())
	homeController := controllers.NewHomeController(servicesProvider.ProvideHomeService())
	categoryController := controllers.NewCategoryController(servicesProvider.ProvideCategoryService())
	workerPublicController := controllers.NewWorkerPublicController(servicesProvider.ProvideWorkerPublicService())
	orderController := controllers.NewOrderController(servicesProvider.ProvideOrderService())
	knowledgeController := controllers.NewKnowledgeController(servicesProvider.ProvideKnowledgeService())
	notificationController := controllers.NewNotificationController(servicesProvider.ProvideNotificationService())
	workerController := controllers.NewWorkerController(servicesProvider.ProvideWorkerService())
	fileController := controllers.NewFileController(servicesProvider.ProvideFileService())
	return &controllerProvider{
		connectionController:   connectionController,
		authController:         authController,
		userController:         userController,
		homeController:         homeController,
		categoryController:     categoryController,
		workerPublicController: workerPublicController,
		orderController:        orderController,
		knowledgeController:    knowledgeController,
		notificationController: notificationController,
		workerController:       workerController,
		fileController:         fileController,
	}
}

func (c *controllerProvider) ProvideConnectionController() controllers.ConnectionController {
	return c.connectionController
}

func (c *controllerProvider) ProvideAuthController() controllers.AuthController {
	return c.authController
}

func (c *controllerProvider) ProvideUserController() controllers.UserController {
	return c.userController
}

func (c *controllerProvider) ProvideHomeController() controllers.HomeController {
	return c.homeController
}

func (c *controllerProvider) ProvideCategoryController() controllers.CategoryController {
	return c.categoryController
}

func (c *controllerProvider) ProvideWorkerPublicController() controllers.WorkerPublicController {
	return c.workerPublicController
}

func (c *controllerProvider) ProvideOrderController() controllers.OrderController {
	return c.orderController
}

func (c *controllerProvider) ProvideKnowledgeController() controllers.KnowledgeController {
	return c.knowledgeController
}

func (c *controllerProvider) ProvideNotificationController() controllers.NotificationController {
	return c.notificationController
}

func (c *controllerProvider) ProvideWorkerController() controllers.WorkerController {
	return c.workerController
}

func (c *controllerProvider) ProvideFileController() controllers.FileController {
	return c.fileController
}
