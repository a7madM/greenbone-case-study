package computers

import (
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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
func UpdateByID(c *fiber.Ctx) error {
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
