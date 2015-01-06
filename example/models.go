package models

import "time"

//go:generate tabler $GOFILE

// User stores user account information.
// @table
type User struct {
	ID      string    `dtype:"uuid" primary:"true"`
	Email   string    `dtype:"varchar(128)"`
	Created time.Time `dtype:"timestamp"`
}

// Profile stores user attributes.
// @table
type Profile struct {
	UserID    string `dtype:"uuid" primary:"true"`
	Attribute string `dtype:"varchar(64)" primary:"true"`
	Value     string `dtype:"varchar(256)"`
}
