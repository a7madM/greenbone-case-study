package computers

import (
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllComputers godoc
// @Summary      Get all computers
// @Description  Retrieves all computer entries
// @Tags         computers
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Computer
// @Failure      500  {object}  map[string]string
// @Router       /api/computers [get]
func GetAll(c *fiber.Ctx) error {
	var computers []models.Computer

	err := database.DB.Find(&computers).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve computers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(computers)
}
