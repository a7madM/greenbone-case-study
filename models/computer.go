package models

import (
	"fmt"
	"greenbone-case-study/utils"

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

	// check abbreviation length
	if len(c.EmployeeAbbreviation) > 3 {
		return fmt.Errorf("employee abbreviation cannot be longer than 3 characters")
	}
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
		return utils.NotifyAdmin(employee, int(count))
	}

	return nil
}
