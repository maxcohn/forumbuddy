package routes

import (
	"database/sql"
	"forumbuddy/repos"
	"forumbuddy/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *appState) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentRepo := repos.CommentRepositorySql{
		DB: app.db,
	}
	// Auth is required on this route

	// Parse the form and validate the values
	r.ParseForm()

	// Validate the form values
	// Make sure the text parameter is present and not empty
	text, err := utils.FormValueStringNonEmpty(r.Form, "text") //TODO: convert to using the payload
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
	_, err = commentRepo.CreateNewComment(curUser.Uid, pid, parentCid, text)

	if err != nil {
		app.render500Page(w)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(pid), http.StatusSeeOther)
}

func (app *appState) commentPageHandler(w http.ResponseWriter, r *http.Request) {
	commentRepo := repos.CommentRepositorySql{
		DB: app.db,
	}
	// Get the comment ID from the router param
	commentId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || commentId <= 0 {
		app.render404Page(w)
		return
	}

	// Get the comment from the DB
	comment, err := commentRepo.GetCommentById(commentId)
	if err != nil {
		app.render404Page(w)
		return
	}

	curUser, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)

	app.templates.ExecuteTemplate(w, "comment.tmpl", map[string]interface{}{
		"Comment":     comment,
		"IsLoggedIn":  isLoggedIn,
		"CurrentUser": curUser,
	})
}
