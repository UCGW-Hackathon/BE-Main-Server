package router

import (
	"log"

	"whatsapp-backend/middleware"
	"whatsapp-backend/provider"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig()
	middleware.ConfigureAuth(config.ProvideJWTConfig().GetSecretKey())
	ConnectionRouter(router, controller)
	err := router.Run(config.ProvideEnvConfig().GetTCPAddress())
	if err != nil {
		log.Fatal("failed to start router:", err)
	}
}
