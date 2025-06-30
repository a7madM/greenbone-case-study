package computers

import (
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

// DeleteByID godoc
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
func DeleteByID(c *fiber.Ctx) error {
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
