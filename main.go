package main

import (
	"fmt"

	"greenbone-case-study/database"
	"greenbone-case-study/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()

	app := fiber.New()
	routes.SetupRoutes(app)

	fmt.Println("Starting server on port 3000...")
	app.Listen(":3000")

}
