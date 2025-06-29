package routes

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Greenbone Case Study API!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("API is running! Current time: %s", time.Now().Format(time.RFC3339)))
	})

	api := app.Group("/api")

	api.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "This is a test endpoint",
		})
	})
}
