package main

import "github.com/jmoiron/sqlx"

var Conn *sqlx.DB = sqlx.MustConnect("postgres", "user=postgres password=password host=127.0.0.1 sslmode=disable")
