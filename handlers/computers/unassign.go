package computers

import (
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

// UnAssign godoc
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
// UnAssign unassigns a computer from an employee by its ID.
// It sets the EmployeeAbbreviation field to an empty string.
// If the computer is not found, it returns a 404 error.
// If the unassignment is successful, it returns the updated computer object.
func UnAssign(c *fiber.Ctx) error {
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
