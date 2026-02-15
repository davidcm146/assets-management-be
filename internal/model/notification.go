package model

import (
	"time"

	"github.com/davidcm146/assets-management-be.git/internal/utils"
)

type NotificationType int

const (
	NotificationTypeLoanSlipOverdue NotificationType = iota + 1
)

type Notification struct {
	ID          int        `json:"id"`
	RecipientID int        `json:"recipient_id"`
	SenderID    *int       `json:"sender_id"`
	Title       string     `json:"title"`
	Type        int        `json:"type"`
	Content     string     `json:"content"`
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (t NotificationType) String() string {
	switch t {
	case NotificationTypeLoanSlipOverdue:
		return "loan_slip_overdue"
	default:
		return "unknown"
	}
}

func ParseType(s string) NotificationType {
	switch s {
	case "loan_slip_overdue":
		return NotificationTypeLoanSlipOverdue
	default:
		return 0
	}
}

func (t NotificationType) MarshalJSON() ([]byte, error) {
	return utils.MarshalEnum(t)
}
