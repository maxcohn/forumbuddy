package models

import "time"

type Post struct {
	Pid       int       `db:"pid" json:"pid"`
	Title     string    `db:"title" json:"title"`
	Body      string    `db:"body" json:"body"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	User     User      `db:"user" json:"user"`
	Comments []Comment `json:"comments"`
}
