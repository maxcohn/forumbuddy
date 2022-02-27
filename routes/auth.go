package routes

import (
	"forumbuddy/models"
	"forumbuddy/utils"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func (app *appState) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	// If the user is already logged in, redirect them to the homepage
	_, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)
	if isLoggedIn {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	app.templates.ExecuteTemplate(w, "login.tmpl", map[string]interface{}{})
}

func (app *appState) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Validate username and password from form
	username, err := utils.FormValueStringNonEmpty(r.Form, "username")
	if err != nil {
		app.templates.ExecuteTemplate(w, "login.tmpl", map[string]interface{}{"Error": "Invalid username or password"})
		return
	}
	username = strings.TrimSpace(username)

	password, err := utils.FormValueStringNonEmpty(r.Form, "password")
	if err != nil {
		app.templates.ExecuteTemplate(w, "login.tmpl", map[string]interface{}{"Error": "Invalid username or password"})
		return
	}

	// Verify the password matches the stored hash and the associated user back
	user, err := models.VerifyUserPassword(app.db, username, password)
	if err != nil {
		app.templates.ExecuteTemplate(w, "login.tmpl", map[string]interface{}{"Error": "Invalid username or password"})
		return
	}

	// Read the session
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.render500Page(w)
		return
	}

	sess.Values["user"] = user

	err = sess.Save(r, w)
	if err != nil {
		app.render500Page(w)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)

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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *appState) signupPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in. If they are, ignore this and redirect them to the homepage
	_, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)
	if isLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the page template
	app.templates.ExecuteTemplate(w, "signup.tmpl", nil)

}

//TODO: change name to signupHandler?
func (app *appState) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in. If they are, ignore this and redirect them to the homepage
	_, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)
	if isLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse form for username and passwords
	r.ParseForm() //TODO: on 400, rerender signup with error message
	username, err := utils.FormValueStringNonEmpty(r.Form, "username")
	if err != nil {
		app.templates.ExecuteTemplate(w, "signup.tmpl", map[string]interface{}{"Error": "Invalid username"})
		return
	}

	password, err := utils.FormValueStringNonEmpty(r.Form, "password")
	if err != nil {
		app.templates.ExecuteTemplate(w, "signup.tmpl", map[string]interface{}{"Error": "Invalid password or passwords do not match"})
		return
	}

	confirmPassword, err := utils.FormValueStringNonEmpty(r.Form, "confirmpassword")
	if err != nil {
		app.templates.ExecuteTemplate(w, "signup.tmpl", map[string]interface{}{"Error": "Invalid password or passwords do not match"})
		return
	}

	// Check if both of their passwords match
	if password != confirmPassword {
		app.templates.ExecuteTemplate(w, "signup.tmpl", map[string]interface{}{"Error": "Invalid password or passwords do not match"})
		return
	}

	// Check if the username exists already
	userExists, err := models.UserExistsByUsername(app.db, username)
	if err != nil {
		app.render500Page(w)
		return
	}

	if userExists {
		app.templates.ExecuteTemplate(w, "signup.tmpl", map[string]interface{}{"Error": "A user with that username already exists"})
		return
	}

	// Hash the password
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		//TODO: Log this error because I don't think the password hash should fail
		app.render500Page(w)
		return
	}

	// Store the new user in the database
	newUser, err := models.CreateNewUser(app.db, username, passwordHash)
	if err != nil {
		app.render500Page(w)
		return
	}

	// Set their session as logged in
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		app.render500Page(w)
		return
	}

	sess.Values["user"] = newUser
	sess.Save(r, w)

	// Redirect them to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
