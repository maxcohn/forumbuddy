package routes

import (
	"forumbuddy/repos"
	"forumbuddy/utils"
	"log"
	"net/http"

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

type LoginPayload struct {
	Username string `validate:"required" form:"username"`
	Password string `validate:"required" form:"password"`
}

func (app *appState) loginUserHandler(w http.ResponseWriter, r *http.Request) {

	userRepo := repos.UserRepositorySql{
		DB: app.db,
	}

	r.ParseForm()

	var loginPayload LoginPayload
	err := utils.DecodeAndValidateForm(&loginPayload, r.Form)
	if err != nil {
		log.Printf("failed to decode and validate: %s", err.Error())
		return
	}

	// Verify the password matches the stored hash and the associated user back
	user, err := userRepo.VerifyUserPassword(loginPayload.Username, loginPayload.Password)
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

type CreateUserPayload struct {
	Username        string `validate:"required" form:"username"`
	Password        string `validate:"required" form:"password"`
	ConfirmPassword string `validate:"required" form:"confirmpassword"`
}

//TODO: change name to signupHandler?
func (app *appState) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in. If they are, ignore this and redirect them to the homepage
	_, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)
	if isLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	userRepo := repos.UserRepositorySql{
		DB: app.db,
	}

	r.ParseForm()

	var createUserPayload CreateUserPayload
	err := utils.DecodeAndValidateForm(&createUserPayload, r.Form)
	if err != nil {
		log.Printf("failed to decode and validate: %s", err.Error())
		return
	}

	// Check if both of their passwords match
	if createUserPayload.Password != createUserPayload.ConfirmPassword {
		app.templates.ExecuteTemplate(w, "signup.tmpl", map[string]interface{}{"Error": "Invalid password or passwords do not match"})
		return
	}

	// Check if the username exists already
	_, err = userRepo.GetUserByUsername(createUserPayload.Username)
	if err == nil {
		app.templates.ExecuteTemplate(w, "signup.tmpl", map[string]interface{}{"Error": "A user with that username already exists"})
		return
	}

	// Hash the password
	passwordHash, err := argon2id.CreateHash(createUserPayload.Password, argon2id.DefaultParams)
	if err != nil {
		//TODO: Log this error because I don't think the password hash should fail
		app.render500Page(w)
		return
	}

	// Store the new user in the database
	newUser, err := userRepo.CreateNewUser(createUserPayload.Username, passwordHash)
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
