package routes

import (
	"forumbuddy/models"
	"forumbuddy/repos"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *appState) userPageHandler(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	userRepo := repos.UserRepositorySql{
		DB: app.db,
	}

	// Check if the route parameter is a string or an int
	if uid, err := strconv.Atoi(chi.URLParam(r, "idOrUsername")); err != nil {
		// If it's a string, we query the user by their username
		username := chi.URLParam(r, "idOrUsername")

		user, err = userRepo.GetUserByUsername(username)

		if err != nil {
			app.render404Page(w)
			return
		}
	} else {
		// If it's an int, we query the user by their uid
		user, err = userRepo.GetUserById(uid)

		if err != nil {
			app.render404Page(w)
			return
		}
	}

	app.templates.ExecuteTemplate(w, "user.tmpl", user)
}
