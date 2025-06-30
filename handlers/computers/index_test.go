package computers

import (
	"encoding/json"
	"greenbone-case-study/database"
	"greenbone-case-study/models"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
