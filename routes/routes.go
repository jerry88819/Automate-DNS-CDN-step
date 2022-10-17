package routes

import (
    "Tencent_backstage_api/handler"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/session"
)

func SetupRoutes(router fiber.Router) {

    sess := session.New()

    tencent := router.Group("/tencent")

	// inject session to request context
	tencent.Use(func(c *fiber.Ctx) error {
		c.Locals("session", sess)
		return c.Next()
	})
    
    // Tencent
    tencent.Post("/test", handler.CreateTencentCdn )
    tencent.Get("/purge", handler.PurgeSingleDomain )


} // SetupNoteRoutes()