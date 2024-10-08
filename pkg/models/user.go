package models

import "time"

type LoginUser struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CompanyID int64  `json:"company_id" binding:"required"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	CompanyID int64     `json:"company_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
