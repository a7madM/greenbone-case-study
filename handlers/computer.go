package handlers

import (
	"fmt"
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateComputer godoc
// @Summary      Create a new computer
// @Description  Creates a computer entry
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        computer  body  models.Computer  true  "Computer"
// @Success      201  {object}  models.Computer
// @Failure      400  {object}  map[string]string
// @Router       /api/computers [post]
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

// GetComputerByID godoc
// @Summary      Get a computer by ID
// @Description  Retrieves a computer entry by its ID
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Computer ID"
// @Success      200  {object}  models.Computer
// @Failure      404  {object}  map[string]string
// @Router       /api/computers/{id} [get]
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

// GetAllComputers godoc
// @Summary      Get all computers
// @Description  Retrieves all computer entries
// @Tags         computers
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Computer
// @Failure      500  {object}  map[string]string
// @Router       /api/computers [get]
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

// DeleteComputerByID godoc
// @Summary      Delete a computer by ID
// @Description  Deletes a computer entry by its ID
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Computer ID"
// @Success      204  {string}  string  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/computers/{id} [delete]
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

// UpdateComputerByID godoc
// @Summary      Update a computer by ID
// @Description  Updates a computer entry by its ID
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Computer ID"
// @Param        computer  body  models.Computer  true  "Computer"
// @Success      200  {object}  models.Computer
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/computers/{id} [put]
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

// AssignComputer godoc
// @Summary      Assign a computer to an employee
// @Description  Assigns a computer to an employee by their abbreviation
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Computer ID"
// @Param        abbr  path  string  true  "Employee Abbreviation"
// @Success      200  {object}  models.Computer
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/computers/{id}/assign/{abbr} [post]
// AssignComputer assigns a computer to an employee by their abbreviation.
// It updates the computer's EmployeeAbbreviation field with the provided abbreviation.
// If the computer is not found, it returns a 404 error.
// If the assignment is successful, it returns the updated computer object.
// If there is an error during the assignment, it returns a 500 error.
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

// GetEmployeeComputers godoc
// @Summary      Get all computers assigned to an employee
// @Description  Retrieves all computers assigned to an employee by their abbreviation
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        abbr  path  string  true  "Employee Abbreviation"
// @Success      200  {array}  models.Computer
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/employees/{abbr}/computers [get]
// GetEmployeeComputers retrieves all computers assigned to an employee by their abbreviation.
// It queries the database for computers where the EmployeeAbbreviation matches the provided abbreviation.
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

// UnAssignComputer godoc
// @Summary      Unassign a computer from an employee
// @Description  Unassigns a computer from an employee by its ID
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Computer ID"
// @Success      200  {object}  models.Computer
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/computers/{id}/unassign [post]
// UnAssignComputer unassigns a computer from an employee by its ID.
// It sets the EmployeeAbbreviation field to an empty string.
// If the computer is not found, it returns a 404 error.
// If the unassignment is successful, it returns the updated computer object.
func UnAssignComputer(c *fiber.Ctx) error {
	id := c.Params("id")

	var computer models.Computer

	if err := database.DB.First(&computer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Can't find Computer with ID " + id})
	}

	computer.EmployeeAbbreviation = ""

	if err := database.DB.Save(&computer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to unassign computer"})
	}
	return c.JSON(computer)
}
