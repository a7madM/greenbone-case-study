package computers

import (
	"encoding/json"
	"greenbone-case-study/database"
	"greenbone-case-study/models"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
