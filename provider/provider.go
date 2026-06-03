package provider

import (
	"situkang/middleware"

	"github.com/gin-gonic/gin"
)

type AppProvider interface {
	ProvideRouter() *gin.Engine
	ProvideConfig() ConfigProvider
	ProvideRepositories() RepositoriesProvider
	ProvideServices() ServicesProvider
	ProvideControllers() ControllerProvider
	ProvideMiddlewares() MiddlewareProvider
}
type appProvider struct {
	ginRouter            *gin.Engine
	configProvider       ConfigProvider
	repositoriesProvider RepositoriesProvider
	servicesProvider     ServicesProvider
	controllerProvider   ControllerProvider
	middlewareProvider   MiddlewareProvider
}

func NewAppProvider() AppProvider {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.SecurityHeaders(), middleware.RateLimitHeaders(100))
	ginRouter.Static("/uploads", "./uploads")
	configProvider := NewConfigProvider()
	repositoriesProvider := NewRepositoriesProvider(configProvider)
	servicesProvider := NewServicesProvider(repositoriesProvider, configProvider)
	controllerProvider := NewControllerProvider(servicesProvider)
	middlewareProvider := NewMiddlewareProvider(servicesProvider)
	// configProvider.ProvideDatabaseConfig().AutoMigrateAll()

	return &appProvider{
		ginRouter:            ginRouter,
		configProvider:       configProvider,
		repositoriesProvider: repositoriesProvider,
		servicesProvider:     servicesProvider,
		controllerProvider:   controllerProvider,
		middlewareProvider:   middlewareProvider,
	}
}
func (a *appProvider) ProvideRouter() *gin.Engine {
	return a.ginRouter
}
func (a *appProvider) ProvideConfig() ConfigProvider {
	return a.configProvider
}

func (a *appProvider) ProvideRepositories() RepositoriesProvider {
	return a.repositoriesProvider
}

func (a *appProvider) ProvideServices() ServicesProvider {
	return a.servicesProvider
}

func (a *appProvider) ProvideControllers() ControllerProvider {
	return a.controllerProvider
}

func (a *appProvider) ProvideMiddlewares() MiddlewareProvider {
	return a.middlewareProvider
}
