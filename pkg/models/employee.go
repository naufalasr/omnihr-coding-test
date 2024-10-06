package models

import "time"

type Status string

const (
	Active     Status = "active"
	NotStarted Status = "not_started"
	Terminated Status = "terminated"
)

type Employee struct {
	ID          int64     `json:"id,omitempty" gorm:"primary_key;autoIncrement"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	ContactInfo string    `json:"contact_info,omitempty"`
	CompanyID   int64     `json:"company_id,omitempty" gorm:"index"`
	Department  string    `json:"department,omitempty" gorm:"index"`
	Position    string    `json:"position,omitempty" gorm:"index"`
	Location    string    `json:"location,omitempty" gorm:"index"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}
