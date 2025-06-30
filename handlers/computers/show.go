package computers

import (
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

// GetByID godoc
// @Summary      Get a computer by ID
// @Description  Retrieves a computer entry by its ID
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Computer ID"
// @Success      200  {object}  models.Computer
// @Failure      404  {object}  map[string]string
// @Router       /api/computers/{id} [get]
func GetByID(c *fiber.Ctx) error {
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
