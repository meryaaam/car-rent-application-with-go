package models

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	Available Status = "available"
	Rented    Status = "rented"
)

type Car struct {
	ID           uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey;default:uuid()"`
	Model        string    `json:"model"`
	Registration string    `json:"registration"`
	Mileage      int       `json:"mileage"`
	Status       Status    `json:"status" gorm:"default:rented"`
	CreatedAt    time.Time `json:"created_at"`
}
