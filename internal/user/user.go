package user

import "time"

type User struct {
	ID              int
	GoogleID        string
	Email           string
	Name            string
	PictureURL      string
	DefaultGliderID int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
