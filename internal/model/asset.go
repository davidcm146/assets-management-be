package model

import "time"

const (
	AssetStatusReturned  = "returned"
	AssetStatusBorrowing = "borrowing"
	AssetStatusAvailable = "available"
)

type Asset struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" binding:"required,gte=3,lte=100"`
	OwnerID      int       `json:"owner_id" binding:"required"`
	Department   string    `json:"department" binding:"required"`
	Position     string    `json:"position" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	Status       string    `json:"status" binding:"oneof=returned borrowing available"`
	SerialNumber int32     `json:"serial_number"`
	Images       []string  `json:"images" binding:"images"`
	CreatedBy    int       `json:"created_by" binding:"required"`
	BorrowedDate time.Time `json:"borrowed_date"`
	ReturnedDate time.Time `json:"returned_date"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}
