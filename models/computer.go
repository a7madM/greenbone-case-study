package models

import (
	"gorm.io/gorm"
)

type Computer struct {
	gorm.Model

	MACAddress           string `json:"mac_address" validate:"required"`
	ComputerName         string `json:"computer_name" validate:"required"`
	IPAddress            string `json:"ip_address" validate:"required"`
	EmployeeAbbreviation string `json:"employee_abbreviation,omitempty"`
	Description          string `json:"description,omitempty"`
}

func (c *Computer) Validate() bool {
	if c.MACAddress == "" || c.ComputerName == "" || c.IPAddress == "" {
		return false
	}
	return true
}
