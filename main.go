package main

import (
	"fmt"
	"forumbuddy/routes"

	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/boj/redistore.v1"
)

/*TODO:
 * user validation via hash (argon2id)
 * Post rating
 * Proper logging
 * caching routes (only if not logged in?)
 *	Maybe add caching to models individually for expensive things like post and comments?
 * Add is_link bool to schema and update models
 * Redis for ratelimiting
 * Migrate error handling to something like this: https://go.dev/blog/error-handling-and-go
 * Add signup link to header
 * replace returns from models with pointers
 */

/* New error handling example:
type RequestError struct {
	Error error
	Message string
	StatusCode int
}

fn (re )

type RequestHandler func(http.ResponseWriter, *http.Request) *RequestError

func (fn RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		// Check error code and handle differently
	}
}
*/
func main() {
	// Setup data being used
	db := sqlx.MustConnect("postgres", "user=postgres password=password host=127.0.0.1 sslmode=disable")
	templates := template.Must(template.New("").ParseGlob("./templates/*.tmpl"))

	secretKey := "thisisexactly32charactersmydudes"

	//TODO: get secret key from env
	//sessionStore := sessions.NewCookieStore([]byte(secretKey))
	sessionStore, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte(secretKey))
	if err != nil {
		panic(err)
	}
	defer sessionStore.Close()

	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//TODO: ping to test conenction
	//redisStatus := redis.Ping(context.Background())
	//log.Fatalf(redisStatus.Err().Error())

	router := routes.NewRouter(db, templates, sessionStore, redis)

	mainRouter := chi.NewRouter()
	mainRouter.Mount("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	//TODO: favicon.ico
	mainRouter.Mount("/favicon.ico", http.FileServer(http.Dir("./static")))
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
