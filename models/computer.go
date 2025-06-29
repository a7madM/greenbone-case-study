package models

import "gorm.io/gorm"

type Computer struct {
	gorm.Model
	MACAddress           string `json:"mac_address" gorm:"not null"`
	ComputerName         string `json:"computer_name" gorm:"not null"`
	IPAddress            string `json:"ip_address" gorm:"not null"`
	EmployeeAbbreviation string `json:"employee_abbreviation,omitempty"`
	Description          string `json:"description,omitempty"`
}
