package model

import (
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/utils"
)

type Status int

const (
	Borrowing Status = iota + 1
	Returned
)

type LoanSlip struct {
	ID           int        `json:"id"`
	Name         string     `json:"name" label:"Tên tài sản"`
	BorrowerName string     `json:"borrower_name" label:"Tên người mượn"`
	Department   string     `json:"department" label:"Phòng ban"`
	Position     string     `json:"position" label:"Chức vụ"`
	Description  string     `json:"description" label:"Mô tả"`
	Status       Status     `json:"status" label:"Trạng thái"`
	SerialNumber string     `json:"serial_number" label:"Số sê ri"`
	Images       []string   `json:"images" label:"Hình ảnh"`
	CreatedBy    int        `json:"created_by"`
	BorrowedDate *time.Time `json:"borrowed_date" label:"Ngày mượn"`
	ReturnedDate *time.Time `json:"returned_date" label:"Ngày trả"`
	UpdatedAt    *time.Time `json:"updated_at"`
	CreatedAt    *time.Time `json:"created_at"`
}

func (s Status) String() string {
	switch s {
	case Borrowing:
		return "borrowing"
	case Returned:
		return "returned"
	default:
		return "Unknown"
	}
}

func StringToStatus(s string) Status {
	switch s {
	case "borrowing":
		return Borrowing
	case "returned":
		return Returned
	default:
		return 0
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	return utils.MarshalEnum(s)
}
