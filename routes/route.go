package routes

import (
	"fmt"
	"greenbone-case-study/handlers/computers"

	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) *fiber.App {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Greenbone Case Study API!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("API is running! Current time: %s", time.Now().Format(time.RFC3339)))
	})

	app.Post("/api/computers", computers.Create)
	app.Get("/api/computers", computers.GetAll)
	app.Get("/api/computers/:id", computers.GetByID)

	app.Delete("/api/computers/:id", computers.DeleteByID)
	app.Put("/api/computers/:id", computers.UpdateByID)
	app.Put("/api/computers/:id/assign/:abbr", computers.Assign)
	app.Get("/api/employees/:abbr/computers", computers.FilterByEmployeeAbbr)
	app.Post("/api/computers/:id/unassign", computers.UnAssign)
	return app
}
