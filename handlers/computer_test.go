package handlers

import (
	"greenbone-case-study/database"
	"testing"

	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	// Use in-memory DB for isolation
	database.DB = database.ConnectInMemoryDB()

	app := fiber.New()
	app.Post("/api/computers", CreateComputer)

	return app
}

func TestCreateComputer(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:66","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestCreateComputerMissingFields(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:66"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCreateComputerInvalidJSON(t *testing.T) {
	app := setupTestApp()
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(`{"mac_address": "invalid_json"`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
func TestCreateComputerDBError(t *testing.T) {
	app := setupTestApp()
	database.DB.Exec("DROP TABLE computers;")

	payload := `{"mac_address":"11:22:33:44:55:66","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}
