package handlers

import (
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

func CreateComputer(c *fiber.Ctx) error {
	var computer models.Computer

	if err := c.BodyParser(&computer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if !computer.Validate() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing required fields",
		})
	}

	if err := database.DB.Create(&computer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create computer",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(computer)
}
