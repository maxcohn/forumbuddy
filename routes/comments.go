package routes

import (
	"database/sql"
	"fmt"
	"forumbuddy/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// TODO: convert to chi like in main.go
)

func (app *appState) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: rqeuire login
	fmt.Println("hit create new commeent")
	r.ParseForm()
	//TODO: validate?

	// Make sure the text parameter is present and not empty

	// Make sure the 'pid' parameter is present and an int greater than zero

	// Check if the 'cid' (parent cid) is valid. If it is, use it, otherwise, make it null

	text := r.Form["text"][0]
	//TODO: better error hadnling
	pid, _ := strconv.Atoi(r.Form["pid"][0])

	var parentCid sql.NullInt64
	if tmpParentCid, err := strconv.Atoi(r.Form["cid"][0]); err != nil {
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

	curUser, isLoggedIn := getUserIdIfLoggedIn(r, app.sessionStore)

	if !isLoggedIn {
		//TODO: 401
		return
	}

	_, err := models.CreateNewComment(app.db, curUser.Uid, pid, parentCid, text)

	if err != nil {
		http.Error(w, "Failed to create the comment", 500)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(pid), 303)
	//TODO: return 201
}

//TODO: change name to signupHandler?
func (app *appState) createUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *appState) commentPageHandler(w http.ResponseWriter, r *http.Request) {
	commentId, _ := strconv.Atoi(chi.URLParam(r, "id"))

	comment, err := models.GetCommentById(app.db, commentId)

	curUser, isLoggedIn := getUserIdIfLoggedIn(r, app.sessionStore)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.templates.ExecuteTemplate(w, "comment.tmpl", map[string]interface{}{
		"Comment":     comment,
		"IsLoggedIn":  isLoggedIn,
		"CurrentUser": curUser,
	})
}
