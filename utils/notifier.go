package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// this url should be injected from environment variables or configuration, but for simplicity,
// it's hardcoded here, but for production I will depend on environment variables.
var MESSAGING_SYSTEM_URL = "http://message_queue:8080/api/notify"

func NotifyAdmin(employeeAbbreviation string, assignedComputersCount int) error {

	type NotificationPayload struct {
		EmployeeAbbreviation string `json:"employeeAbbreviation"`
		Level                string `json:"level"`
		Message              string `json:"message"`
	}

	payload := NotificationPayload{
		EmployeeAbbreviation: employeeAbbreviation,
		Level:                "warning",
		Message:              fmt.Sprintf("The employee with abbreviation %s has %d computers assigned.", employeeAbbreviation, assignedComputersCount),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(MESSAGING_SYSTEM_URL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to notify admin: %s", resp.Status)
	}
	return nil
}
