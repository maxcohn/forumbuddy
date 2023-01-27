package routes

import (
	"encoding/gob"
	"forumbuddy/models"
	"forumbuddy/repos"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type appState struct {
	db           *sqlx.DB
	templates    *template.Template
	sessionStore sessions.Store
	redis        *redis.Client
}

func (app *appState) requireLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.sessionStore.Get(r, "session")
		if err != nil {
			// If there was an error reading the sessions, we can't confirm they're logged in
			//TODO: log that we failed to read the session
			app.render500Page(w)
			return
		}

		if sess.Values["user"] == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			// If they do have a session, but there is no 'uid' value, that means they're not logged in
			//TODO: maybe redirect to /login?
			return //TODO: report 401
		}

		//ctx := context.WithValue(r.Context(), "username", sess.Values["username"])
		//ctx = context.WithValue(ctx, "uid", sess.Values["uid"])
		//next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)

	})
}

// Session by pointer or not? TODO: maybe a better way to do this?
func getUserIfLoggedIn(r *http.Request, session sessions.Store) (models.User, bool) { //TODO: return entire user struct
	sess, err := session.Get(r, "session")
	if err != nil {
		// If we failed to read the session, we can't confirm they're logged in
		//TODO: figure out the best structure for this
		return models.User{}, false
	}

	if sess.Values["user"] == nil {
		// If the uid is nil, they're not logged in
		return models.User{}, false
	}

	return sess.Values["user"].(models.User), true
}

// Another idea for handling appstate could be to have `AppState` in another package and have it imported by
// other routes, but then typedefed so we can make methods on them that won't clash
/*
type test appState

func (t *test) postPageHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	postId, _ := strconv.Atoi(routeVars["id"])

	post, _ := models.GetPostAndCommentsById(t.db, postId)

	t.templates.ExecuteTemplate(w, "post.tmpl", post)
}
*/

func (app *appState) render404Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

		// Check if there was a 404
		log.Println(r.Response.StatusCode)

	})
}

//TODO: left off with user login For now, maybe we should ignore argon2id and hashing and just use a map for checking
//TODO: maybe switch from gorilla mux to chi?
func NewRouter(db *sqlx.DB, templates *template.Template, sessionStore sessions.Store, redis *redis.Client) *chi.Mux {
	app := appState{
		db:           db,
		templates:    templates,
		sessionStore: sessionStore,
		redis:        redis,
	}

	// Register the User model with gob so we can save it in the session
	gob.Register(models.User{})

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	//router.Use(app.render404Middleware) //TODO: 404 handler, or just make a wrapper function

	//TODO: move elsewhere
	// Set up 404 handler
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("uhto"))
		w.WriteHeader(http.StatusNotFound)
	})

	//router.Use(middleware.RealIP)
	//TODO: test this so we can use the isloggedin middleware  - loggedInRouter := router.NewRoute().Subrouter()

	//TODO: maybe move the route creating to each file? e.g. app.authRouter, app.commentRouter?
	// Page rendering routes
	/*router.Route("/newpost", func(r chi.Router) {
		r.Use(app.requireLoggedInMiddleware)
		r.Get("/", app.newPostPageHandler)
	})*/
	router.Get("/", AppHandler(app.indexHandler).ServeHTTP)
	router.Get("/newpost", app.requireLoggedInMiddleware(http.HandlerFunc(app.newPostPageHandler)).ServeHTTP)
	router.Get("/post/{id:[0-9]+}", AppHandler(app.postPageHandler).ServeHTTP)
	router.Get("/comment/{id:[0-9]+}", app.commentPageHandler)
	router.Get("/user/{idOrUsername}", app.userPageHandler)

	// Creation routes
	router.Post("/post", app.requireLoggedInMiddleware(AppHandler(app.createPostHandler)).ServeHTTP)             //TODO: require loggedin
	router.Post("/comment", app.requireLoggedInMiddleware(http.HandlerFunc(app.createCommentHandler)).ServeHTTP) //TODO: require loggedin
	//TODO: route for user signup

	// Authentication related routes
	router.Get("/login", app.loginPageHandler)
	router.Post("/login", AppHandler(app.loginUserHandler).ServeHTTP)
	router.Get("/logout", app.requireLoggedInMiddleware(http.HandlerFunc(app.logoutUserHandler)).ServeHTTP)
	router.Get("/signup", app.signupPageHandler)
	router.Post("/signup", app.createUserHandler) //TODO: change name?

	return router
}

func (app *appState) indexHandler(w http.ResponseWriter, r *http.Request) AppError {
	postRepo := repos.PostRepositorySql{
		DB: app.db,
	}
	// Get the 10 most recent posts
	posts, err := postRepo.GetRecentPosts(10)

	if err != nil {
		return InternalAppError{}
	}

	// Check if the user is logged in to show login status
	curUser, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)

	app.templates.ExecuteTemplate(w, "index.tmpl", map[string]interface{}{
		"Posts":       posts,
		"IsLoggedIn":  isLoggedIn,
		"CurrentUser": curUser,
	})

	log.Println(map[string]interface{}{
		"IsLoggedIn":  isLoggedIn,
		"CurrentUser": curUser,
	})
	return nil
}

func (app *appState) render500Page(w http.ResponseWriter) {
	w.WriteHeader(500)
	app.templates.ExecuteTemplate(w, "500.tmpl", nil)
}

func (app *appState) render404Page(w http.ResponseWriter) {
	w.WriteHeader(404)
	app.templates.ExecuteTemplate(w, "404.tmpl", nil)
}

/*
func (app *appState) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement rate limiting
	})
}*/

//type AppHandler func(http.Handler) AppError
type AppHandler func(http.ResponseWriter, *http.Request) AppError

func (handler AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appErr := handler(w, r)

	if appErr == nil {
		return
	}

	appErr.Render(w)

	// Check error type
	/*
		switch appErr.(type) {
		case NotFoundAppError:
			return
		case InternalAppError:
			return
		default:
			//TODO: unknown error
		}*/

	/*
		func (fn RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
			err := fn(w, r)
			if err != nil {
				// Check error code and handle differently
			}
		}
	*/
}

//TODO: make middleware to handle different errors
//TODO: switch to using app errors
type AppError interface {
	error
	UserError() string
	Render(http.ResponseWriter)
}

//TODO: general error?

type NotFoundAppError struct{}

func (err NotFoundAppError) UserError() string {
	return ""
}

func (err NotFoundAppError) Render(w http.ResponseWriter) {
	w.WriteHeader(404) //TODO: add template rendering
	w.Write([]byte("oi bruv, no page"))
}

func (err NotFoundAppError) Error() string {
	return ""
}

type InternalAppError struct{}

func (err InternalAppError) UserError() string {
	return ""
}

func (err InternalAppError) Render(w http.ResponseWriter) {
	w.WriteHeader(500) //TODO: add template rendering
}

func (err InternalAppError) Error() string {
	return ""
}
