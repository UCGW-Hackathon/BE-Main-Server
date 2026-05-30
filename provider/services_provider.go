package provider

import "whatsapp-backend/services"

type ServicesProvider interface {
	ProvideConnectionService() services.ConnectionService
	ProvideAuthService() services.AuthService
	ProvideUserService() services.UserService
	ProvideHomeService() services.HomeService
	ProvideCategoryService() services.CategoryService
	ProvideWorkerPublicService() services.WorkerPublicService
	ProvideOrderService() services.OrderService
	ProvideKnowledgeService() services.KnowledgeService
	ProvideNotificationService() services.NotificationService
	ProvideWorkerService() services.WorkerService
}

type servicesProvider struct {
	connectionService   services.ConnectionService
	authService         services.AuthService
	userService         services.UserService
	homeService         services.HomeService
	categoryService     services.CategoryService
	workerPublicService services.WorkerPublicService
	orderService        services.OrderService
	knowledgeService    services.KnowledgeService
	notificationService services.NotificationService
	workerService       services.WorkerService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	db := configProvider.ProvideDatabaseConfig().GetInstance()
	connectionService := services.NewConnectionService(repoProvider.ProvideConnectionRepository())
	authService := services.NewAuthService(db, configProvider.ProvideJWTConfig(), configProvider.ProvideEnvConfig())
	userService := services.NewUserService(db)
	categoryService := services.NewCategoryService(db)
	workerPublicService := services.NewWorkerPublicService(db)
	homeService := services.NewHomeService(db, categoryService, workerPublicService)
	orderService := services.NewOrderService(db)
	knowledgeService := services.NewKnowledgeService(db)
	notificationService := services.NewNotificationService(db)
	workerService := services.NewWorkerService(db)
	return &servicesProvider{
		connectionService:   connectionService,
		authService:         authService,
		userService:         userService,
		homeService:         homeService,
		categoryService:     categoryService,
		workerPublicService: workerPublicService,
		orderService:        orderService,
		knowledgeService:    knowledgeService,
		notificationService: notificationService,
		workerService:       workerService,
	}
}

func (s *servicesProvider) ProvideConnectionService() services.ConnectionService {
	return s.connectionService
}

func (s *servicesProvider) ProvideAuthService() services.AuthService {
	return s.authService
}

func (s *servicesProvider) ProvideUserService() services.UserService {
	return s.userService
}

func (s *servicesProvider) ProvideHomeService() services.HomeService {
	return s.homeService
}

func (s *servicesProvider) ProvideCategoryService() services.CategoryService {
	return s.categoryService
}

func (s *servicesProvider) ProvideWorkerPublicService() services.WorkerPublicService {
	return s.workerPublicService
}

func (s *servicesProvider) ProvideOrderService() services.OrderService {
	return s.orderService
}

func (s *servicesProvider) ProvideKnowledgeService() services.KnowledgeService {
	return s.knowledgeService
}

func (s *servicesProvider) ProvideNotificationService() services.NotificationService {
	return s.notificationService
}

func (s *servicesProvider) ProvideWorkerService() services.WorkerService {
	return s.workerService
}
