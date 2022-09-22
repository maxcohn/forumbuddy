package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Uid       int       `db:"uid" json:"uid"`
	Username  string    `db:"username" json:"username"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Post struct {
	Pid       int       `db:"pid" json:"pid"`
	Title     string    `db:"title" json:"title"`
	Body      string    `db:"body" json:"body"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	User     User      `db:"user" json:"user"`
	Comments []Comment `json:"comments"`
}

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

func GetRecentPosts(db *sqlx.DB, limit int) ([]Post, error) {
	var posts []Post
	err := db.Select(&posts, `
		SELECT
			p.pid,
			p.title,
			p.created_at,
			u.username AS "user.username"
		FROM posts as p, users AS u
		WHERE p.uid = u.uid
		ORDER by p.created_at desc 
		LIMIT $1
	`, limit)

	if err != nil {
		fmt.Println(err.Error())
		//TODO: add logging
		return nil, err
	}

	return posts, nil
}

func GetUserByUsername(db *sqlx.DB, username string) (*User, error) {
	user := new(User)
	err := db.Get(user, `
		SELECT uid, username, created_at
		FROM users
		WHERE username = $1
	`, username)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserById(db *sqlx.DB, id int) (*User, error) {
	user := new(User)
	err := db.Get(user, `
		SELECT uid, username, created_at
		FROM users
		WHERE uid = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetCommentById(db *sqlx.DB, id int) (*Comment, error) {
	comment := new(Comment)

	err := db.Get(comment, `
		SELECT
			c.cid,
			c.pid,
			c.body,
			c.parent,
			c.created_at,
			u.uid AS "user.uid",
			u.username AS "user.username"
		FROM comments AS c, users AS u
		WHERE c.uid = u.uid
			AND c.cid = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

//TODO: get comment and children by id

func GetPostAndCommentsById(db *sqlx.DB, id int) (*Post, error) {
	post := new(Post)
	var comments []Comment

	// Query the current post
	err := db.Get(post, `
		SELECT
			p.pid,
			p.title,
			p.body,
			p.created_at,
			u.uid AS "user.uid",
			u.username AS "user.username"
		FROM posts AS p, users AS u
		WHERE p.uid = u.uid
			AND p.pid = $1
		
	`, id)

	if err != nil {
		//TODO: log error
		return nil, err
	}

	// Query all comments on the post
	err = db.Select(&comments, `
		SELECT
			c.cid,
			c.body,
			c.parent,
			c.pid,
			u.uid AS "user.uid",
			u.username AS "user.username"
		FROM comments AS c, users AS u
		WHERE c.uid = u.uid
			AND c.pid = $1
		ORDER BY
			CASE WHEN parent IS NULL THEN 0
			ELSE parent
		END ASC
	`, id)

	if err != nil {
		return nil, err
	}

	// Create an empty slice of the comments at the root of the tree. These are pointers since we're going to be updating the slices as we go
	rootComments := make([]*Comment, 0)

	// Create an empty mapping from cids to comments. These are points because we're going to be modifying them in our loop
	commentMap := make(map[int]*Comment)

	for i, comment := range comments {
		var curComment = &comments[i]

		curComment.Children = make([]*Comment, 0)
		commentMap[comment.Cid] = curComment

		if !comment.Parent.Valid {
			// If there is not parent comment, this is at the root
			rootComments = append(rootComments, curComment)
		} else {
			// If there is a parent comment
			parent := commentMap[int(comment.Parent.Int64)]
			parent.Children = append(parent.Children, curComment)
		}
	}

	// Convert all root comments to their values
	for _, c := range rootComments {
		post.Comments = append(post.Comments, *c)
	}

	return post, nil
}

func CreateNewPost(db *sqlx.DB, uid int, title, body string) (int, error) {
	// Insert the post into the DB and get that new post's ID
	var newPostId int
	err := db.Get(&newPostId, `
		INSERT INTO posts
			(uid, title, body)
		VALUES
			($1, $2, $3)
		RETURNING pid
	`, uid, title, body)

	if err != nil {
		fmt.Println(err.Error())
		//TODO: log this?
		return 0, err
	}

	return newPostId, nil
}

func CreateNewComment(db *sqlx.DB, uid, pid int, parent sql.NullInt64, body string) (int, error) {
	// Insert the comemnt into the DB and get that new comment's ID
	var newCommentId int
	err := db.Get(&newCommentId, `
		INSERT INTO comments
			(uid, pid, parent, body)
		VALUES
			($1, $2, $3, $4)
		RETURNING cid
	`, uid, pid, parent, body)

	if err != nil {
		fmt.Println(err.Error())
		//TODO: log this?
		return 0, err
	}

	return newCommentId, nil
}

func UserExistsByUsername(db *sqlx.DB, username string) (bool, error) {
	var tossAway int
	// Check if the user exists
	err := db.Get(&tossAway, "SELECT 1 FROM users WHERE username = $1", username)

	// If the reported error is that there are no rows, that means the user doesn't exist
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		// If there was an error returned otherwise, that's an actual DB error
		log.Println("Error checking for user: ", err.Error()) //TODO: better logging
		return false, err
	}

	// If we had any results, the user does exist
	return true, nil
}

//TODO: shoudl this take teh raw password isntead so we can force it to hash it?
func CreateNewUser(db *sqlx.DB, username, passwordHash string) (*User, error) {
	newUser := new(User)
	err := db.Get(newUser, `
		WITH user_ins AS (
			INSERT INTO users
				(username)
			VALUES
				($1)
			RETURNING uid, username
		),
		hash_ins AS (
			INSERT INTO user_hashes
				(uid, password_hash)
			VALUES
				((SELECT uid FROM user_ins), $2)
			RETURNING uid
		)
		SELECT uid, username FROM user_ins
	`, username, passwordHash)

	// Check if the user insert failed
	if err != nil {
		//TODO: log error
		return nil, err
	}

	return newUser, nil
}

//TODO: shoudl this take teh raw password isntead so we can force it to hash it?
func VerifyUserPassword(db *sqlx.DB, username, password string) (*User, error) {
	// Get the hash from the DB for this user
	var passwordHash string
	err := db.Get(&passwordHash, `
		SELECT uh.password_hash
		FROM users AS u, user_hashes AS uh
		WHERE u.uid = uh.uid
			AND u.username = $1
	`, username)

	if err == sql.ErrNoRows { //TODO: different response for no match?
		return nil, err
	} else if err != nil {
		return nil, err
	}

	// Verify password matches the stored password hash
	match, err := argon2id.ComparePasswordAndHash(password, passwordHash)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, errors.New("Password hash didn't match")
	}

	// Now that we know the hashes match, query the user
	return GetUserByUsername(db, username)
}
