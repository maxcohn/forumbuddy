package models

import "time"

type User struct {
	Uid       int       `db:"uid" json:"uid"`
	Username  string    `db:"username" json:"username"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
