package main

import (
	"situkang/provider"
	"situkang/router"
)

func main() {
	appProvider := provider.NewAppProvider()
	router.RunRouter(appProvider)
}
