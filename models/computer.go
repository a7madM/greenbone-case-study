package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// TODO:
// 1. Add validation for MAC address format
// 2. Add validation for IP address format

type Computer struct {
	ID                   uint   `json:"id" gorm:"primaryKey"`
	MACAddress           string `json:"mac_address" gorm:"not null"`
	ComputerName         string `json:"computer_name" gorm:"not null"`
	IPAddress            string `json:"ip_address" gorm:"not null"`
	EmployeeAbbreviation string `json:"employee_abbreviation,omitempty"`
	Description          string `json:"description,omitempty"`
}

func (c *Computer) Validate() bool {
	if c.MACAddress == "" || c.ComputerName == "" || c.IPAddress == "" {
		return false
	}
	return true
}

func (c *Computer) BeforeSave(tx *gorm.DB) (err error) {
	if !c.Validate() {
		return gorm.ErrInvalidData
	}
	var count int64

	// Check for duplicate MAC address
	tx.Model(&Computer{}).Where("mac_address = ? AND id != ?", c.MACAddress, c.ID).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	tx.Model(&Computer{}).Where("ip_address = ? AND id != ?", c.IPAddress, c.ID).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	return nil
}

func (c *Computer) AfterSave(tx *gorm.DB) (err error) {
	employee := c.EmployeeAbbreviation
	if employee == "" {
		return nil
	}
	var count int64
	tx.Model(&Computer{}).Where("employee_abbreviation = ?", employee).Count(&count)

	if count >= 3 {
		return NotifyAdmin(employee)
	}

	return nil
}

// this url should be injected from environment variables or configuration, but for simplicity, it's hardcoded here.
var MESSAGING_SYSTEM_URL = "http://message_queue:8080/api/notify"

func NotifyAdmin(employeeAbbreviation string) error {

	type NotificationPayload struct {
		EmployeeAbbreviation string `json:"employeeAbbreviation"`
		Level                string `json:"level"`
		Message              string `json:"message"`
	}

	payload := NotificationPayload{
		EmployeeAbbreviation: employeeAbbreviation,
		Level:                "warning",
		Message:              "some message",
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
