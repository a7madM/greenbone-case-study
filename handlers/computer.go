package handlers

import (
	"fmt"
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
				"error": fmt.Sprintf("MAC Address %s or IP Address %s already exists", computer.MACAddress, computer.IPAddress)})
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
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

func UpdateComputerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var computer models.Computer

	if err := database.DB.First(&computer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "computer not found"})
	}

	var input models.Computer
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	computer.MACAddress = input.MACAddress
	computer.ComputerName = input.ComputerName
	computer.IPAddress = input.IPAddress
	computer.EmployeeAbbreviation = input.EmployeeAbbreviation
	computer.Description = input.Description

	err := database.DB.Save(&computer).Error

	if err != nil {

		if err == gorm.ErrInvalidData {

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid data provided. Please ensure all required fields are filled correctly.",
			})
		}
	}
	return c.JSON(computer)
}

func AssignComputer(c *fiber.Ctx) error {
	id := c.Params("id")
	abbr := c.Params("abbr")
	fmt.Println("Assigning computer with ID:", id, "to employee with abbreviation:", abbr)

	var computer models.Computer

	if err := database.DB.First(&computer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Can't find Computer with ID " + id})
	}

	computer.EmployeeAbbreviation = abbr

	if err := database.DB.Save(&computer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to assign computer"})
	}
	return c.JSON(computer)
}

func GetEmployeeComputers(c *fiber.Ctx) error {
	abbr := c.Params("abbr")
	fmt.Println("Fetching computers for employee with abbreviation:", abbr)

	var computers []models.Computer

	err := database.DB.Where("employee_abbreviation = ?", abbr).Find(&computers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch employee computers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(computers)
}
