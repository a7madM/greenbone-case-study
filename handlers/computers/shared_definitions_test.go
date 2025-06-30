package computers

import (
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp() *fiber.App {
	database.DB = database.ConnectInMemoryDB()
	database.DB.Exec("DELETE FROM computers")
	database.DB.AutoMigrate(&models.Computer{})

	app := fiber.New()
	SetupComputerRoutes(app)
	return app
}

func SetupComputerRoutes(app *fiber.App) {
	app.Post("/api/computers", Create)
	app.Get("/api/computers", GetAll)
	app.Get("/api/computers/:id", GetByID)
	app.Delete("/api/computers/:id", DeleteByID)
	app.Put("/api/computers/:id", UpdateByID)
	app.Put("/api/computers/:id/assign/:abbr", Assign)
	app.Get("/api/employees/:abbr/computers", FilterByEmployeeAbbr)
	app.Post("/api/computers/:id/unassign", UnAssign)
}

func createComputer(name, macAddress, ipAddr, employeeAbbreviation string) models.Computer {
	computer := models.Computer{
		MACAddress:           macAddress,
		ComputerName:         name,
		IPAddress:            ipAddr,
		EmployeeAbbreviation: employeeAbbreviation,
	}
	database.DB.Create(&computer)
	return computer
}
