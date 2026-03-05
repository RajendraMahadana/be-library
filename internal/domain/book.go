package domain

import "time"

type Book struct {
	ID        uint
	Title     string
	Author    string
	ISBN      string
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
