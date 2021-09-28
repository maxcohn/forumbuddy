package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Uid          int       `db:"uid" json:"uid"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
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

	User     User      `db:"user" json:"user"`
	Children []Comment `json:"children"`
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := Conn.Get(&user, `
		SELECT uid, username, created_at
		FROM users
		WHERE username = $1
	`, username)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserById(id int) (User, error) {
	var user User
	err := Conn.Get(&user, `
		SELECT uid, username, created_at
		FROM users
		WHERE uid = $1
	`, id)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetCommentById(id int) (Comment, error) {
	var comment Comment

	err := Conn.Get(&comment, `
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

func GetPostAndCommentsById(id int) (Post, error) {
	var post Post
	var comments []Comment

	// Query the current post
	err := Conn.Get(&post, `
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
		fmt.Println(err)
		return post, err
	}

	// Query all comments on the post
	err = Conn.Select(&comments, `
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

		curComment.Children = make([]Comment, 0)
		commentMap[comment.Cid] = curComment

		if !comment.Parent.Valid {
			fmt.Printf("%v\n", *curComment)
			// If there is not parent comment, this is at the root
			rootComments = append(rootComments, curComment)
		} else {
			// If there is a parent comment
			parent := commentMap[int(comment.Parent.Int64)]
			parent.Children = append(parent.Children, *curComment)
		}
	}

	// Convert all root comments to their values
	for _, c := range rootComments {
		post.Comments = append(post.Comments, *c)
	}

	return post, nil
}

func main() {
	/*TODO:
	 * user validation via hash (argon2id)
	 * Middleware for auth above
	 * cookie setup middleware?
	 * Gorrilla mux routing
	 * Create new post/user/comment
	 * Post rating
	 * Session expiration
	 * templates
	 */

	templates := template.Must(template.New("post").ParseGlob("./templates/*.tmpl"))

	router := mux.NewRouter()
	router.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)
		postId, _ := strconv.Atoi(routeVars["id"])

		post, _ := GetPostAndCommentsById(postId)

		templates.ExecuteTemplate(w, "post.tmpl", post)

	})

	router.HandleFunc("/comment/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)
		commentId, _ := strconv.Atoi(routeVars["id"])

		comment, err := GetCommentById(commentId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		templates.ExecuteTemplate(w, "comment.tmpl", comment)
	})

	router.HandleFunc("/user/{idOrUsername}", func(w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)

		var user User

		if uid, err := strconv.Atoi(routeVars["idOrUsername"]); err != nil {
			username := routeVars["idOrUsername"]

			user, err = GetUserByUsername(username)

			if err != nil {
				http.Error(w, "Unable to find user", http.StatusNotFound)
				return
			}
		} else {
			user, err = GetUserById(uid)

			if err != nil {
				http.Error(w, "Unable to find user", http.StatusNotFound)
				return
			}
		}

		templates.ExecuteTemplate(w, "user.tmpl", user)
	})

	srv := &http.Server{
		Addr: "127.0.0.1:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	srv.ListenAndServe()
}
