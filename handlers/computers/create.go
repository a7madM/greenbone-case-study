package computers

import (
	"fmt"
	"greenbone-case-study/database"
	"greenbone-case-study/models"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary      Create a new computer
// @Description  Creates a computer entry
// @Tags         computers
// @Accept       json
// @Produce      json
// @Param        computer  body  models.Computer  true  "Computer"
// @Success      201  {object}  models.Computer
// @Failure      400  {object}  map[string]string
// @Router       /api/computers [post]
func Create(c *fiber.Ctx) error {
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
