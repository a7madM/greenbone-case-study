package computers

import (
	"fmt"
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

// Assign godoc
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
// Assign assigns a computer to an employee by their abbreviation.
// It updates the computer's EmployeeAbbreviation field with the provided abbreviation.
// If the computer is not found, it returns a 404 error.
// If the assignment is successful, it returns the updated computer object.
// If there is an error during the assignment, it returns a 500 error.
func Assign(c *fiber.Ctx) error {
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
