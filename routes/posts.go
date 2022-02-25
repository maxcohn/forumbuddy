package routes

import (
	"fmt"
	"forumbuddy/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" // TODO: convert to chi like in main.go
)

func (app *appState) newPostPageHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: require login - if not logged in, redirect to /login

	app.templates.ExecuteTemplate(w, "newpost.tmpl", nil)
}

func (app *appState) postPageHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	postId, _ := strconv.Atoi(routeVars["id"])

	_, isLoggedIn := getUserIdIfLoggedIn(r, app.sessionStore)
	fmt.Println("IsLoggedIn? ", isLoggedIn)

	post, _ := models.GetPostAndCommentsById(app.db, postId)

	app.templates.ExecuteTemplate(w, "post.tmpl", map[string]interface{}{
		"Post":       post,
		"IsLoggedIn": isLoggedIn,
	})
}

func (app *appState) createPostHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: rqeuire login
	fmt.Println("hit create new post")
	r.ParseForm()
	//TODO: validate?
	title := r.Form["title"][0]
	text := r.Form["text"][0]
	uid, isLoggedIn := getUserIdIfLoggedIn(r, app.sessionStore)

	if !isLoggedIn {
		//TODO: 401
		return
	}

	newPostId, err := models.CreateNewPost(app.db, uid, title, text)

	if err != nil {
		http.Error(w, "Failed to create the post", 500)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(newPostId), 303)
	//TODO: return 201
}
