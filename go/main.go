package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
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

	var users []User
	err := Conn.Select(&users, "SELECT * FROM users")

	if err != nil {
		panic(err)
	}

	postTemplate := template.Must(template.New("post").ParseFiles("./templates/post.tmpl"))

	fmt.Println(postTemplate.Name())
	post, _ := GetPostAndCommentsById(1)

	v, err := json.MarshalIndent(post, "", "    ")

	fmt.Printf("%v\n", string(v))
	fmt.Printf("%v\n", post)

	err = postTemplate.ExecuteTemplate(os.Stdout, "post.tmpl", post)

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)
		postId, _ := strconv.Atoi(routeVars["id"])

		post, _ := GetPostAndCommentsById(postId)

		postTemplate := template.Must(template.New("post").ParseFiles("./templates/post.tmpl"))

		postTemplate.ExecuteTemplate(w, "post.tmpl", post)

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
