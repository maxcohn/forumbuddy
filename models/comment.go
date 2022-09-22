package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	Cid       int           `db:"cid" json:"cid"`
	Pid       int           `db:"pid" json:"pid"`
	Uid       int           `db:"uid" json:"-"`
	Body      string        `db:"body" json:"body"`
	Parent    sql.NullInt64 `db:"parent" json:"-"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`

	User     User       `db:"user" json:"user"`
	Children []*Comment `json:"children"`
}
