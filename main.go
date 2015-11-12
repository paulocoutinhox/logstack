package main

import (
	"github.com/prsolucoes/logstack/app"
	"github.com/prsolucoes/logstack/controllers"
)

func main() {
	app.Server = app.NewWebServer()
	app.Server.LoadConfiguration()
	app.Server.CreateDataSource()
	app.Server.CreateBasicRoutes()

	{
		controller := controllers.HomeController{}
		controller.Register()
	}

	{
		controller := controllers.LogController{}
		controller.Register()
	}

	{
		controller := controllers.APIController{}
		controller.Register()
	}

	app.Server.Start()
}
