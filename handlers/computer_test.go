package handlers

import (
	"encoding/json"
	"fmt"
	"greenbone-case-study/database"
	"greenbone-case-study/models"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	// Use in-memory DB for isolation
	// Clear any existing data and create fresh tables

	database.DB = database.ConnectInMemoryDB()
	database.DB.AutoMigrate(&models.Computer{})
	app := fiber.New()
	app.Post("/api/computers", CreateComputer)
	app.Get("/api/computers/:id", GetComputerByID)
	app.Get("/api/computers", GetAllComputers)

	return app
}

func TestCreateComputer(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:60","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestCreateComputerMissingFields(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:60"}`
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

func TestCreateComputerDuplicateMAC(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:61","computer_name":"TestPC","ip_address":"192.168.1.4"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	// Try to create the same computer again
	req = httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, `{"error":"computer with this MAC or IP address already exists"}`, string(body))
	assert.Equal(t, 409, resp.StatusCode)
}

func TestCreateComputerDuplicateIP(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:68","computer_name":"TestPC","ip_address":"192.168.1.5"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	payload = `{"mac_address":"11:22:33:44:55:69","computer_name":"TestPC","ip_address":"192.168.1.5"}`
	req = httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"computer with this MAC or IP address already exists"}`, string(body))
	assert.Equal(t, 409, resp.StatusCode)
}

func TestGetAllComputersEmpty(t *testing.T) {
	app := setupTestApp()
	database.DB.Exec("DELETE FROM computers WHERE 1=1") // Clear any existing records

	req := httptest.NewRequest("GET", "/api/computers", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)

	var computers []models.Computer
	err = json.Unmarshal(body, &computers)
	assert.Nil(t, err)
	fmt.Println("Computers:", computers)
}

func TestGetAllComputers(t *testing.T) {
	app := setupTestApp()
	// Clear any existing records
	database.DB.Exec("DELETE FROM computers")

	payload := `{"mac_address":"11:22:33:44:55:66","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	fmt.Println("Error:", err)

	req = httptest.NewRequest("GET", "/api/computers", nil)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"mac_address":"11:22:33:44:55:66"`)
	assert.Contains(t, string(body), `"computer_name":"TestPC"`)
	assert.Contains(t, string(body), `"ip_address":"192.168.1.2"`)
}

func TestCreateComputerDBError(t *testing.T) {
	app := setupTestApp()
	database.DB.Exec("DROP TABLE computers;")

	payload := `{"mac_address":"11:22:33:44:55:60","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetComputerByID(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:68","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)

	computer := models.Computer{}
	err = json.Unmarshal(body, &computer)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/computers/%d", computer.ID), nil)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
func TestGetComputerByIDNotFound(t *testing.T) {
	app := setupTestApp()
	fakeID := 9999
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/computers/%d", fakeID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"computer not found"}`, string(body))
}
