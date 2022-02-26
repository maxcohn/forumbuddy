package main

import (
	"fmt"
	"forumbuddy/routes"

	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
* caching routes (only if not logged in?)
*	Maybe add caching to models individually for expensive things like post and comments?
* fix templates using base template
* Add is_link bool to schema and update models
 */
func main() {

	// Setup data being used
	db := sqlx.MustConnect("postgres", "user=postgres password=password host=127.0.0.1 sslmode=disable")
	templates := template.Must(template.New("").ParseGlob("./templates/*.tmpl"))

	secretKey := "thisisexactly32charactersmydudes"

	//TODO: swap this out for redis
	//TODO: get secret key from env
	sessionStore := sessions.NewCookieStore([]byte(secretKey))

	router := routes.NewRouter(db, templates, sessionStore)

	mainRouter := chi.NewRouter()
	mainRouter.Mount("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	//TODO: is this hacky or acceptable? I'd rather it be at the root for standards purposes instead of in static, but having it in static would be nicer

	mainRouter.Mount("/", router)

	srv := &http.Server{
		Addr: "127.0.0.1:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mainRouter, // Pass our instance of gorilla/mux in.
	}

	fmt.Println("Starting server")

	srv.ListenAndServe()
}
