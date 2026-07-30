package models

import "time"

type Session struct {
	ID        string
	AdminID   string
	Token     string
	ExpiresAt time.Time
}
