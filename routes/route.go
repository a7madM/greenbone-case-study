package routes

import (
	"fmt"
	"greenbone-case-study/handlers"
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
	app.Post("/api/computers", handlers.CreateComputer)
	app.Get("/api/computers/:id", handlers.GetComputerByID)
	app.Get("/api/computers", handlers.GetAllComputers)
	app.Delete("/api/computers/:id", handlers.DeleteComputerByID)
	app.Put("/api/computers/:id", handlers.UpdateComputerByID)
	app.Put("/api/computers/:id/assign/:abbr", handlers.AssignComputer)
	app.Get("/api/employees/:abbr/computers", handlers.GetEmployeeComputers)
}
