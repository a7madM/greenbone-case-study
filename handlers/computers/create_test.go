package computers

import (
	"fmt"
	"greenbone-case-study/database"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	createComputer("TestPC1", "11:22:33:44:55:90", "192.168.1.2", "EM1")
	createComputer("TestPC2", "11:22:33:44:55:91", "192.168.1.3", "EM1")

	payload := `{"mac_address":"11:22:33:44:55:94","computer_name":"TestPC","ip_address":"192.168.1.6","employee_abbreviation":"EM1"}`
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
