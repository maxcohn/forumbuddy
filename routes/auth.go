package routes

import (
	"log"
	"net/http"
)

func (app *appState) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		//TODO: figure out how to handle this case (since it shouldn't happen, I think, maybe just log and 500)
		return
	}

	// If the user is already logged in, redirect them to the homepage
	if sess.Values["uid"] != nil {
		//TODO: remove? sess.Values["uid"] = nil
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	app.templates.ExecuteTemplate(w, "login.tmpl", map[string]interface{}{
		//"User": sess.Values["username"],
	})
}

func (app *appState) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form["username"][0] //TODO: better error handling/validation here
	//password := r.Form["password"]

	//TODO: add password checking. No password checking at the moment for development

	//TODO: move this to models this
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
	log.Println("username?: ", username)

	//TODO: hadnle error
	sess, err := app.sessionStore.Get(r, "session")
	if err != nil {
		//TODO: handle error
		log.Println("Error reading session")
		return
	}

	sess.Values["uid"] = uid
	//sess.Values["username"] = username

	sess.Save(r, w)

	// Validate username and password

	// Hash and compare passowrd

	// Redirect to home page
	http.Redirect(w, r, "/", 303)

}

func (app *appState) logoutUserHandler(w http.ResponseWriter, r *http.Request) {

}