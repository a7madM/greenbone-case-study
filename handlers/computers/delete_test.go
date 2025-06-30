package computers

import (
	"encoding/json"
	"fmt"
	"greenbone-case-study/database"
	"greenbone-case-study/models"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
