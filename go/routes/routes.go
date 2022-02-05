package routes

import (
	"forumbuddy/models"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
			return
		}

		if sess.Values["uid"] == nil {
			// If they do have a session, but there is no 'uid' value, that means they're not logged in
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

func NewRouter(db *sqlx.DB, templates *template.Template, sessionStore sessions.Store) *mux.Router {
	app := appState{
		db:           db,
		templates:    templates,
		sessionStore: sessionStore,
	}

	router := mux.NewRouter()
	//TODO: test this so we can use the isloggedin middleware  - loggedInRouter := router.NewRoute().Subrouter()

	// Page rendering routes
	router.HandleFunc("/", app.indexHandler).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}", app.postPageHandler).Methods("GET")
	router.HandleFunc("/comment/{id:[0-9]+}", app.commentPageHandler).Methods("GET")
	router.HandleFunc("/user/{idOrUsername}", app.userPageHandler).Methods("GET")

	// User related routes
	router.HandleFunc("/login", app.loginPageHandler).Methods("GET")

	router.HandleFunc("/login", app.loginUserHandler).Methods("POST")

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

func (app *appState) commentPageHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	commentId, _ := strconv.Atoi(routeVars["id"])

	comment, err := models.GetCommentById(app.db, commentId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.templates.ExecuteTemplate(w, "comment.tmpl", comment)
}

func (app *appState) postPageHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	postId, _ := strconv.Atoi(routeVars["id"])

	post, _ := models.GetPostAndCommentsById(app.db, postId)

	app.templates.ExecuteTemplate(w, "post.tmpl", post)
}

func (app *appState) userPageHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)

	var user models.User

	if uid, err := strconv.Atoi(routeVars["idOrUsername"]); err != nil {
		username := routeVars["idOrUsername"]

		user, err = models.GetUserByUsername(app.db, username)

		if err != nil {
			http.Error(w, "Unable to find user", http.StatusNotFound)
			return
		}
	} else {
		user, err = models.GetUserById(app.db, uid)

		if err != nil {
			http.Error(w, "Unable to find user", http.StatusNotFound)
			return
		}
	}

	app.templates.ExecuteTemplate(w, "user.tmpl", user)
}

func (app *appState) createPostHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *appState) createCommentHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *appState) createUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *appState) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		//TODO:
		return
	}

	// If the user is already logged in, redirect them to the homepage
	if sess.Values["uid"] != nil {
		sess.Values["uid"] = nil
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	app.templates.ExecuteTemplate(w, "login.tmpl", nil)
}

func (app *appState) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form["username"][0] //TODO: better error handling/validation here
	//password := r.Form["password"]

	//TODO: add password checking. No password checking at the moment for development

	var uid int
	//err := app.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&isMatch)
	err := app.db.QueryRowx(`SELECT uid from users as u where username = $1`, username).Scan(&uid)

	if err != nil {
		// If there is an error, there were no rows
		log.Println("Failed to login")
		return
	}

	if err != nil {
		log.Println("Error: ", err.Error())
	}

	log.Println("uid?: ", uid)

	//TODO: hadnle error
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		//TODO: handle error
		log.Println("Error reading session")
		return
	}

	sess.Values["uid"] = uid

	sess.Save(r, w)

	// Validate username and password

	// Hash and compare passowrd
}

func (app *appState) logoutUserHandler(w http.ResponseWriter, r *http.Request) {

}
