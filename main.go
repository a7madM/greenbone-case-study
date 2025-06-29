package main

import (
	"fmt"

	"greenbone-case-study/database"
	"greenbone-case-study/routes"

	"github.com/gofiber/fiber/v2"

	_ "greenbone-case-study/docs"

	"github.com/gofiber/swagger"
)

func main() {
	database.ConnectDB()

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.SetupRoutes(app)

	fmt.Println("Starting server on port 3000...")
	app.Listen(":3000")

}
