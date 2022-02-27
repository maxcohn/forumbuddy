package routes

import (
	"forumbuddy/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *appState) userPageHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Check if the route parameter is a string or an int
	if uid, err := strconv.Atoi(chi.URLParam(r, "idOrUsername")); err != nil {
		// If it's a string, we query the user by their username
		username := chi.URLParam(r, "idOrUsername")

		user, err = models.GetUserByUsername(app.db, username)

		if err != nil {
			app.render404Page(w)
			return
		}
	} else {
		// If it's an int, we query the user by their uid
		user, err = models.GetUserById(app.db, uid)

		if err != nil {
			app.render404Page(w)
			return
		}
	}

	app.templates.ExecuteTemplate(w, "user.tmpl", user)
}

//TODO: change name to signupHandler?
func (app *appState) createUserHandler(w http.ResponseWriter, r *http.Request) {
	//TODO:
}
