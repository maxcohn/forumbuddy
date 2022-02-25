package routes

import (
	"forumbuddy/models"
	"github.com/gorilla/mux" // TODO: convert to chi like in main.go
	"net/http"
	"strconv"
)

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
