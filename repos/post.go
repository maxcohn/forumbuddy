package repos

import (
	"forumbuddy/models"
)

type PostRepository interface {
	GetPostAndCommentsById(id int) (*models.Post, error)
	CreateNewPost(uid int, title, body string) (int, error)
	GetRecentPosts(limit int) ([]models.Post, error)
}
