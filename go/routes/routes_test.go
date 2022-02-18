package routes_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //TODO: possibly swap out for sqlite???
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	// setup DB connection

	db = sqlx.MustConnect("postgres", "user=postgres password=password host=127.0.0.1 sslmode=disable")

	// Fill DB with mockdata TODO: - could use sqlite instead

	//routes.NewRouter(db, templates, sessionStore)

	m.Run()

}

//TODO: check proper stauses on all router
//TODO: check login pass and fail
//TODO: check create user
//TODO: check crud post and comments
