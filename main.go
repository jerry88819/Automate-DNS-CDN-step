package main

import (

	"Tencent_backstage_api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Start a new fiber app
	app := fiber.New()

	// Setup the router	api := app.Group("/library", logger.New())
	api := app.Group("", logger.New())

	// Setup the Node Routes
	routes.SetupRoutes(api)

	// Listen on PORT 3000
	app.Listen(":3000")

} // main()