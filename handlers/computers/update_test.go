package computers

import (
	"encoding/json"
	"fmt"
	"greenbone-case-study/models"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	computer := createComputer("TestPC", "11:22:33:44:55:90", "192.168.1.2", "EM1")
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
	computer := createComputer("TestPC", "11:22:33:44:55:90", "192.168.1.2", "EM1")
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/computers/%d", computer.ID), strings.NewReader(`{"mac_address": "invalid_json"`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"invalid body"}`, string(body))
}
