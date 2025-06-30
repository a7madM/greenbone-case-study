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

func TestGetByID(t *testing.T) {
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
func TestGetByIDNotFound(t *testing.T) {
	app := setupTestApp()
	fakeID := 9999
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/computers/%d", fakeID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"computer not found"}`, string(body))
}
