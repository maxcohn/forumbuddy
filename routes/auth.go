package routes

import (
	"forumbuddy/models"
	"forumbuddy/utils"
	"log"
	"net/http"
	"strings"
)

func (app *appState) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.render500Page(w)
	}

	// If the user is already logged in, redirect them to the homepage
	if sess.Values["user"] != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	app.templates.ExecuteTemplate(w, "login.tmpl", map[string]interface{}{})
}

func (app *appState) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Validate username from form
	username, err := utils.FormValueStringNonEmpty(r.Form, "username")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	username = strings.TrimSpace(username)
	//password := r.Form["password"]

	//TODO: add password checking. No password checking at the moment for development

	//TODO: move this to models this
	var user models.User
	//err := app.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&isMatch)
	err = app.db.Get(&user, `SELECT uid, username from users as u where username = $1`, username)

	if err != nil {
		// If there is an error, there were no rows
		log.Println("Failed to login")
		//TODO: return status code
		return
	}

	if err != nil {
		log.Println("Error: ", err.Error())
	}

	// Read the session
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.render500Page(w)
		return
	}

	sess.Values["user"] = user
	//sess.Values["username"] = username

	err = sess.Save(r, w)
	if err != nil {
		app.render500Page(w)
		return
	}

	// Validate username and password

	// Hash and compare passowrd

	// Redirect to home page
	http.Redirect(w, r, "/", 303)

}

func (app *appState) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	// Auth required

	// Get session
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.render500Page(w)
		return
	}

	// Clear the user from the session and save it
	sess.Values["user"] = nil

	err = sess.Save(r, w)
	if err != nil {
		app.render500Page(w)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", 303)
}
