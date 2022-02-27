package models_test

import (
	"forumbuddy/models"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	// setup DB connection
	db = sqlx.MustConnect("postgres", "user=postgres password=password host=127.0.0.1 sslmode=disable")

	// Fill DB with mockdata TODO:

	m.Run()

}

func TestGetRecentPosts(t *testing.T) {
	posts, err := models.GetRecentPosts(db, 10)

	if err != nil {
		t.Errorf("Error querying recent posts: %s", err.Error())
	}

	if len(posts) > 10 {
		t.Errorf("Returned too many posts")
	}
}

func TestGetUserByUsername(t *testing.T) {
	//TODO: remove uid check and just make sure a user is gotten?
	testCases := []models.User{
		{Uid: 1, Username: "foo-bar"},
		{Uid: 2, Username: "pg"},
		{Uid: 3, Username: "ken"},
		{Uid: 4, Username: "dmr"},
	}

	for _, tc := range testCases {
		user, err := models.GetUserByUsername(db, tc.Username)
		if err != nil {
			t.Errorf("Error getting user from DB: %s", err.Error())
		}

		if tc.Uid != user.Uid {
			t.Errorf("Test case uid and result uid do not match. tc: %d, result: %d", tc.Uid, user.Uid)
		}

		if tc.Username != user.Username {
			t.Errorf("Test case username and result username do not match. tc: %s, result: %s", tc.Username, user.Username)
		}
	}
}

func TestGetUserById(t *testing.T) {
	//TODO: move this outside with the username case?
	//TODO: remove uid check and just make sure a user is gotten?
	testCases := []models.User{
		{Uid: 1, Username: "foo-bar"},
		{Uid: 2, Username: "pg"},
		{Uid: 3, Username: "ken"},
		{Uid: 4, Username: "dmr"},
	}

	for _, tc := range testCases {
		user, err := models.GetUserById(db, tc.Uid)
		if err != nil {
			t.Errorf("Error getting user from DB: %s", err.Error())
		}

		if tc.Uid != user.Uid {
			t.Errorf("Test case uid and result uid do not match. tc: %d, result: %d", tc.Uid, user.Uid)
		}

		if tc.Username != user.Username {
			t.Errorf("Test case username and result username do not match. tc: %s, result: %s", tc.Username, user.Username)
		}
	}
}

//TODO: get comment by id - need more test data in sql
//TODO: get post and comments - need more test data, will need to do a reflect.deepcopy
