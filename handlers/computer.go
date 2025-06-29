package handlers

import (
	"fmt"
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
			"error": "name, ip_address, MAC address are all required fields",
		})
	}

	err := database.DB.Create(&computer).Error
	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "computer with this MAC or IP address already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create computer",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(computer)
}

func GetComputerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var computer models.Computer

	err := database.DB.First(&computer, id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "computer not found",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(computer)
}

func GetAllComputers(c *fiber.Ctx) error {
	fmt.Println("Retrieving all computers")
	var computers []models.Computer

	err := database.DB.Find(&computers).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve computers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(computers)
}

func DeleteComputerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var computer models.Computer

	err := database.DB.First(&computer, id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "computer not found",
			})
		}
	}

	err = database.DB.Delete(&computer).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete computer",
		})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}
