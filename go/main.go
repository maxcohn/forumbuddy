package main

import (
	"fmt"
	"forumbuddy/routes"

	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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
func main() {

	// Setup data being used
	db := sqlx.MustConnect("postgres", "user=postgres password=password host=127.0.0.1 sslmode=disable")
	templates := template.Must(template.New("post").ParseGlob("./templates/*.tmpl"))

	secretKey := "thisisexactly32charactersmydudes"

	//TODO: swap this out for redis
	//TODO: get secret key from env
	sessionStore := sessions.NewCookieStore([]byte(secretKey))

	router := routes.NewRouter(db, templates, sessionStore)

	srv := &http.Server{
		Addr: "127.0.0.1:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	fmt.Println("Starting server")

	srv.ListenAndServe()
}
