package routes

import (
	"database/sql"
	"forumbuddy/models"
	"forumbuddy/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *appState) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Auth is required on this route

	// Parse the form and validate the values
	r.ParseForm()

	// Validate the form values
	// Make sure the text parameter is present and not empty
	text, err := utils.FormValueStringNonEmpty(r.Form, "text")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Make sure the 'pid' parameter is present and an int greater than zero
	pid, err := utils.FormValueIntGtZero(r.Form, "pid")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Check if the 'cid' (parent cid) is valid. If it is, use it, otherwise, make it null
	var parentCid sql.NullInt64
	if tmpParentCid, err := strconv.Atoi(r.Form.Get("cid")); err != nil {
		parentCid = sql.NullInt64{
			Valid: false,
			Int64: 0,
		}
	} else {
		parentCid = sql.NullInt64{
			Valid: true,
			Int64: int64(tmpParentCid),
		}
	}

	// Get the current user
	curUser, _ := getUserIfLoggedIn(r, app.sessionStore)

	// Insert the new comment into the DB
	_, err = models.CreateNewComment(app.db, curUser.Uid, pid, parentCid, text)

	if err != nil {
		http.Error(w, "Failed to create the comment", 500)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(pid), 303)
}

func (app *appState) commentPageHandler(w http.ResponseWriter, r *http.Request) {
	// Get the comment ID from the router param
	commentId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || commentId <= 0 {
		w.WriteHeader(404)
		app.templates.ExecuteTemplate(w, "404.tmpl", nil)
		return
	}

	// Get the comment from the DB
	comment, err := models.GetCommentById(app.db, commentId)
	if err != nil {
		w.WriteHeader(404)
		app.templates.ExecuteTemplate(w, "404.tmpl", nil)
		return
	}

	curUser, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)

	app.templates.ExecuteTemplate(w, "comment.tmpl", map[string]interface{}{
		"Comment":     comment,
		"IsLoggedIn":  isLoggedIn,
		"CurrentUser": curUser,
	})
}
