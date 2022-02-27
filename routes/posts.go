package routes

import (
	"forumbuddy/models"
	"forumbuddy/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Render the 'newpost' page template
func (app *appState) newPostPageHandler(w http.ResponseWriter, r *http.Request) {
	app.templates.ExecuteTemplate(w, "newpost.tmpl", nil)
}

// Render the 'post' page for the requested post
func (app *appState) postPageHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post id from the URL
	postId, err := strconv.Atoi(chi.URLParam(r, "id"))

	// If the post id wasn't specified or is invalid, return a 404
	if err != nil { //TODO: make a wrapper for this that takes a string for a context message
		app.render404Page(w)
		return
	}

	// Check if the user is logged in to show login status
	curUser, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)

	// Get the post and its comments from the database
	post, err := models.GetPostAndCommentsById(app.db, postId)

	// If we couldn't get the post, it doesn't exist, so return a 404
	if err != nil {
		app.render404Page(w)
		return
	}

	// Render the 'post' page
	app.templates.ExecuteTemplate(w, "post.tmpl", map[string]interface{}{
		"Post":        post,
		"IsLoggedIn":  isLoggedIn,
		"CurrentUser": curUser,
	})
}

// Create a new post based on the given parameters
func (app *appState) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Auth is required for this route

	// Parse the form and validate the values
	r.ParseForm()

	// Get the values from the form and validate that they are non-empty strings
	title, err := utils.FormValueStringNonEmpty(r.Form, "title")
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	text, err := utils.FormValueStringNonEmpty(r.Form, "text")
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	// Get the current user (we already verified they're logged in via middleware)
	curUser, _ := getUserIfLoggedIn(r, app.sessionStore)

	// Create the new post in the database and get its id
	newPostId, err := models.CreateNewPost(app.db, curUser.Uid, title, text)

	if err != nil {
		app.render500Page(w)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(newPostId), http.StatusSeeOther)
}
