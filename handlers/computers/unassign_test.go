package computers

import (
	"encoding/json"
	"fmt"
	"greenbone-case-study/models"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnAssignComputer(t *testing.T) {
	app := setupTestApp()
	computer := createComputer("TestPC", "11:22:33:44:55:90", "192.168.1.2", "EM1")
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/computers/%d/unassign", computer.ID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var updatedComputer models.Computer
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &updatedComputer)
	assert.Nil(t, err)
	assert.Equal(t, "", updatedComputer.EmployeeAbbreviation)
}

func TestUnAssignComputerNotFound(t *testing.T) {
	app := setupTestApp()
	fakeID := 9999
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/computers/%d/unassign", fakeID), nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"error":"Can't find Computer with ID 9999"}`, string(body))
}
