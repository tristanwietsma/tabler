package models

import "time"

//go:generate tabler $GOFILE

// User stores user account information.
// @table
type User struct {
	ID      string    `type:uuid,primary:true`
	Email   string    `type:varchar(128)`
	Created time.Time `type:timestamp`
}

// Profile stores user attributes.
// @table
type Profile struct {
	UserID    string `type:uuid,primary:true`
	Attribute string `type:varchar(64),primary:true`
	Value     string `type:varchar(256)`
}
