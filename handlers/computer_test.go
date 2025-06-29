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

	database.DB = database.ConnectInMemoryDB()
	database.DB.Exec("DELETE FROM computers")
	database.DB.AutoMigrate(&models.Computer{})
	app := fiber.New()
	app.Post("/api/computers", CreateComputer)
	app.Get("/api/computers/:id", GetComputerByID)
	app.Get("/api/computers", GetAllComputers)
	app.Delete("/api/computers/:id", DeleteComputerByID)
	app.Put("/api/computers/:id", UpdateComputerByID)
	app.Put("/api/computers/:id/assign/:abbr", AssignComputer)
	app.Get("/api/employees/:abbr/computers", GetEmployeeComputers)
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

func TestCreateComputerWhenEmployeeAbbreviationIsOverAssigned(t *testing.T) { // 3 computers or more are assigned to the same employee
	app := setupTestApp()

	createComputer("TestPC1", "11:22:33:44:55:90", "192.168.1.2", "EMP1")
	createComputer("TestPC2", "11:22:33:44:55:91", "192.168.1.3", "EMP1")

	payload := `{"mac_address":"11:22:33:44:55:94","computer_name":"TestPC","ip_address":"192.168.1.6","employee_abbreviation":"EMP1"}`
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

	assert.Equal(t, fmt.Sprintf(`{"error":"MAC Address %s or IP Address %s already exists"}`, "11:22:33:44:55:61", "192.168.1.4"), string(body))
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
	assert.Equal(t, fmt.Sprintf(`{"error":"MAC Address %s or IP Address %s already exists"}`, "11:22:33:44:55:69", "192.168.1.5"), string(body))
	assert.Equal(t, 409, resp.StatusCode)
}

func TestGetAllComputersEmpty(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/api/computers", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)

	var computers []models.Computer
	err = json.Unmarshal(body, &computers)
	assert.Nil(t, err)
}

func TestGetAllComputersError(t *testing.T) {
	app := setupTestApp()
	database.DB.Exec("DROP TABLE computers;")
	req := httptest.NewRequest("GET", "/api/computers", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetAllComputers(t *testing.T) {
	app := setupTestApp()

	payload := `{"mac_address":"11:22:33:44:55:66","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)

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
func TestDeleteComputer(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:70","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("POST", "/api/computers", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	body, _ := io.ReadAll(resp.Body)
	computer := models.Computer{}
	json.Unmarshal(body, &computer)
	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/computers/%d", computer.ID), nil)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestDeleteComputerNotFound(t *testing.T) {
	app := setupTestApp()
	fakeID := 9999
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/computers/%d", fakeID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"computer not found"}`, string(body))
}
func TestDeleteComputerDBError(t *testing.T) {
	app := setupTestApp()
	database.DB.Exec("DROP TABLE computers;")

	fakeID := 9999
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/computers/%d", fakeID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"failed to delete computer"}`, string(body))
}
func TestUpdateComputerByID(t *testing.T) {
	app := setupTestApp()
	payload := `{"mac_address":"11:22:33:44:55:80","computer_name":"TestPC","ip_address":"192.168.1.2"}`
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

	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/computers/%d", computer.ID), strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
func TestUpdateComputerByIDNotFound(t *testing.T) {
	app := setupTestApp()
	fakeID := 9999
	payload := `{"mac_address":"11:22:33:44:55:90","computer_name":"TestPC","ip_address":"192.168.1.2"}`
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/computers/%d", fakeID), strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"computer not found"}`, string(body))
}

func TestUpdateComputerByIDWithEmptyIP(t *testing.T) {
	app := setupTestApp()
	computer := createComputer("TestPC", "11:22:33:44:55:90", "192.168.1.2", "EMP1")
	payload := `{"mac_address":"11:22:33:44:55:90","computer_name":"TestPC","ip_address":""}`
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/computers/%d", computer.ID), strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"Invalid data provided. Please ensure all required fields are filled correctly."}`, string(body))
}

func TestUpdateComputerByIDWithInvalidJSON(t *testing.T) {
	app := setupTestApp()
	computer := createComputer("TestPC", "11:22:33:44:55:90", "192.168.1.2", "EMP1")
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/computers/%d", computer.ID), strings.NewReader(`{"mac_address": "invalid_json"`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"invalid body"}`, string(body))
}

func TestAssignComputer(t *testing.T) {
	app := setupTestApp()
	computer := createComputer("TestPC", "11:22:33:44:55:90", "192.168.1.2", "EMP1")
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/computers/%d/assign/EMP2", computer.ID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var updatedComputer models.Computer
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &updatedComputer)
	assert.Nil(t, err)
	assert.Equal(t, "EMP2", updatedComputer.EmployeeAbbreviation)
}
func TestAssignComputerNotFound(t *testing.T) {
	app := setupTestApp()
	fakeID := 9999
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/computers/%d/assign/EMP2", fakeID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"Can't find Computer with ID 9999"}`, string(body))
}

func TestGetEmployeeComputers(t *testing.T) {
	app := setupTestApp()
	createComputer("TestPC1", "11:22:33:44:55:91", "192.168.1.3", "EMP1")
	createComputer("TestPC2", "11:22:33:44:55:92", "192.168.1.4", "EMP1")

	req := httptest.NewRequest("GET", "/api/employees/EMP1/computers", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var computers []models.Computer
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &computers)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(computers))
}

func TestGetEmployeeComputerError(t *testing.T) {
	app := setupTestApp()
	database.DB.Exec("DROP TABLE computers;")

	req := httptest.NewRequest("GET", "/api/employees/EMP1/computers", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"failed to fetch employee computers"}`, string(body))
}

func createComputer(name, macAddress, ipAddr, employeeAbbreviation string) models.Computer {
	computer := models.Computer{
		MACAddress:           macAddress,
		ComputerName:         name,
		IPAddress:            ipAddr,
		EmployeeAbbreviation: employeeAbbreviation,
	}
	database.DB.Create(&computer)
	return computer
}
