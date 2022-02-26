package routes

import (
	"forumbuddy/models"
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
	if err != nil {
		http.Error(w, "Post not found", 404)
		//TODO: make a template for the 404
		return
	}

	// Check if the user is logged in to show login status
	_, isLoggedIn := getUserIdIfLoggedIn(r, app.sessionStore)

	// Get the post and its comments from the database
	post, err := models.GetPostAndCommentsById(app.db, postId)
	if err != nil {
		http.Error(w, "There was an error getting the requests post and comments", 500)
		return
	}

	app.templates.ExecuteTemplate(w, "post.tmpl", map[string]interface{}{
		"Post":        post,
		"CurrentUser": isLoggedIn,
	})
}

// Create a new post based on the given parameters
func (app *appState) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Auth is required for this route
	r.ParseForm()
	//TODO: validate?
	title := r.Form["title"][0]
	text := r.Form["text"][0]
	uid, _ := getUserIdIfLoggedIn(r, app.sessionStore)

	newPostId, err := models.CreateNewPost(app.db, uid, title, text)

	if err != nil {
		http.Error(w, "Failed to create the post", 500)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(newPostId), 303)
}
