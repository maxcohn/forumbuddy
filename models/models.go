package models

import (
	"database/sql"
	"fmt"
	"time"

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

func GetUserByUsername(db *sqlx.DB, username string) (User, error) {
	var user User
	err := db.Get(&user, `
		SELECT uid, username, created_at
		FROM users
		WHERE username = $1
	`, username)

	if err != nil {
		//TODO: should I be returning a pointer here so I can pass null
		return user, err
	}

	return user, nil
}

func GetUserById(db *sqlx.DB, id int) (User, error) {
	var user User
	err := db.Get(&user, `
		SELECT uid, username, created_at
		FROM users
		WHERE uid = $1
	`, id)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetCommentById(db *sqlx.DB, id int) (Comment, error) {
	var comment Comment

	err := db.Get(&comment, `
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
		return comment, err
	}

	return comment, nil
}

//TODO: get comment and children by id

func GetPostAndCommentsById(db *sqlx.DB, id int) (Post, error) {
	var post Post
	var comments []Comment

	// Query the current post
	err := db.Get(&post, `
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
		return post, err
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
		return post, err
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

//TODO: func CreateNewUser
