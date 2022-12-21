package router

import "github.com/gofiber/fiber/v2"

func SetRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"msg": "Service running.",
		})
	})

	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Set("Version", "v1")
		return c.Next()
	})

	loadRecordsRoutes(v1)
}
