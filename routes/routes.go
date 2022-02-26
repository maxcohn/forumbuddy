package routes

import (
	"forumbuddy/models"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5" // TODO: convert to chi like in main.go
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type appState struct {
	db           *sqlx.DB
	templates    *template.Template
	sessionStore sessions.Store
}

func (app *appState) requireLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.sessionStore.Get(r, "session")
		if err != nil {
			// If there was an error reading the sessions, we can't confirm they're logged in
			//TODO: log that we failed to read the session
			w.WriteHeader(500)
			w.Write([]byte("There was an error on our side"))
			return
		}

		if sess.Values["uid"] == nil {
			http.Redirect(w, r, "/login", 303)
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
func getUserIdIfLoggedIn(r *http.Request, session sessions.Store) (int, bool) {
	sess, err := session.Get(r, "session")
	if err != nil {
		// If we failed to read the session, we can't confirm they're logged in
		//TODO: figure out the best structure for this
		return 0, false
	}

	if sess.Values["uid"] == nil {
		// If the uid is nil, they're not logged in
		return 0, false
	}

	return sess.Values["uid"].(int), true
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

//TODO: left off with user login For now, maybe we should ignore argon2id and hashing and just use a map for checking
//TODO: maybe switch from gorilla mux to chi?
func NewRouter(db *sqlx.DB, templates *template.Template, sessionStore sessions.Store) *chi.Mux {
	app := appState{
		db:           db,
		templates:    templates,
		sessionStore: sessionStore,
	}

	router := chi.NewRouter()
	//TODO: test this so we can use the isloggedin middleware  - loggedInRouter := router.NewRoute().Subrouter()

	// Page rendering routes
	router.Get("/", app.indexHandler)
	router.Get("/newpost", app.newPostPageHandler)
	router.Get("/post/{id:[0-9]+}", app.postPageHandler)
	router.Get("/comment/{id:[0-9]+}", app.commentPageHandler)
	router.Get("/user/{idOrUsername}", app.userPageHandler)

	router.Post("/post", app.createPostHandler)       //TODO: require loggedin
	router.Post("/comment", app.createCommentHandler) //TODO: require loggedin

	// User related routes
	router.Get("/login", app.loginPageHandler)

	router.Post("/login", app.loginUserHandler)

	return router
}

func (app *appState) indexHandler(w http.ResponseWriter, r *http.Request) {

	// Get the 10 most recent posts
	posts, err := models.GetRecentPosts(app.db, 10)

	if err != nil {
		http.Error(w, "Failed to get most recent posts", 400)
		return
	}

	app.templates.ExecuteTemplate(w, "index.tmpl", map[string]interface{}{
		"Posts": posts,
	})
}
