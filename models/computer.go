package models

import (
	"gorm.io/gorm"
)

type Computer struct {
	gorm.Model
	ID                   uint   `json:"id" gorm:"primaryKey"`
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

func (c *Computer) BeforeCreate(tx *gorm.DB) (err error) {
	var count int64

	tx.Model(&Computer{}).Where("mac_address = ?", c.MACAddress).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	tx.Model(&Computer{}).Where("ip_address = ?", c.IPAddress).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	return nil
}
