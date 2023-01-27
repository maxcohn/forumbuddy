package routes

import (
	"forumbuddy/repos"
	"forumbuddy/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Render the 'newpost' page template
func (app *appState) newPostPageHandler(w http.ResponseWriter, r *http.Request) {
	app.templates.ExecuteTemplate(w, "newpost.tmpl", nil)
}

// Render the 'post' page for the requested post
func (app *appState) postPageHandler(w http.ResponseWriter, r *http.Request) AppError {
	postRepo := repos.PostRepositorySql{
		DB: app.db,
	}

	// Get the post id from the URL
	postId, err := strconv.Atoi(chi.URLParam(r, "id"))

	// If the post id wasn't specified or is invalid, return a 404
	if err != nil {
		return NotFoundAppError{}
	}

	// Check if the user is logged in to show login status
	curUser, isLoggedIn := getUserIfLoggedIn(r, app.sessionStore)

	// Get the post and its comments from the database
	post, err := postRepo.GetPostAndCommentsById(postId)

	// If we couldn't get the post, it doesn't exist, so return a 404
	if err != nil {
		//app.render404Page(w)
		return NotFoundAppError{}
	}

	// Render the 'post' page
	app.templates.ExecuteTemplate(w, "post.tmpl", map[string]interface{}{
		"Post":        post,
		"IsLoggedIn":  isLoggedIn,
		"CurrentUser": curUser,
	})

	return nil
}

type CreatePostPayload struct {
	Title string `validate:"required" form:"title"`
	Text  string `validate:"required" form:"text"`
}

// Create a new post based on the given parameters
func (app *appState) createPostHandler(w http.ResponseWriter, r *http.Request) AppError {
	// Auth is required for this route
	postRepo := repos.PostRepositorySql{
		DB: app.db,
	}

	// Parse the form and validate the values
	r.ParseForm()

	var createPostPayload CreatePostPayload
	err := utils.DecodeAndValidateForm(&createPostPayload, r.Form)
	if err != nil {
		log.Printf("failed to decode and validate: %s", err.Error())
		return InternalAppError{} //TODO: return validation error or generic 400
	}

	// Get the current user (we already verified they're logged in via middleware)
	curUser, _ := getUserIfLoggedIn(r, app.sessionStore)

	// Create the new post in the database and get its id
	newPostId, err := postRepo.CreateNewPost(curUser.Uid, createPostPayload.Title, createPostPayload.Text)

	if err != nil {
		return InternalAppError{}
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(newPostId), http.StatusSeeOther)

	return nil
}
