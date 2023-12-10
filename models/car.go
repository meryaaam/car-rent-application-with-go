package models

import (
	"time"
)

type Status string

const (
	Available Status = "available"
	Rented    Status = "rented"
)

type Car struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Model        string    `json:"model"`
	Registration string    `json:"registration"`
	Mileage      int       `json:"mileage"`
	Status       Status    `json:"status" gorm:"default:rented"`
	CreatedAt    time.Time `json:"created_at"`
}
