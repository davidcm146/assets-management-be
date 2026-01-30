package model

import "time"

type Role int

const (
	Admin Role = iota + 1 // Admin = 1
	IT                    // IT = 2
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" binding:"required,gte=3,lte=10"`
	Password  string    `json:"password" binding:"required,gte=6,lte=20"`
	Role      Role      `json:"role" binding:"required,oneof=1 2"`
	CreatedAt time.Time `json:"created_at"`
}
