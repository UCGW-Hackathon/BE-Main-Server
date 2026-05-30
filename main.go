package main

import (
	"whatsapp-backend/provider"
	"whatsapp-backend/router"
)

func main() {
	appProvider := provider.NewAppProvider()
	router.RunRouter(appProvider)
}
