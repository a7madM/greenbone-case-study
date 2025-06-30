package computers

import (
	"fmt"
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

// FilterByEmployeeAbbr godoc
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
// FilterByEmployeeAbbr retrieves all computers assigned to an employee by their abbreviation.
// It queries the database for computers where the EmployeeAbbreviation matches the provided abbreviation.
func FilterByEmployeeAbbr(c *fiber.Ctx) error {
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
