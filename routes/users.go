package routes

import (
	"forumbuddy/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// TODO: convert to chi like in main.go
)

func (app *appState) userPageHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if uid, err := strconv.Atoi(chi.URLParam(r, "idOrUsername")); err != nil {
		username := chi.URLParam(r, "idOrUsername")

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
