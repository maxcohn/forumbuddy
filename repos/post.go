package repos

import (
	"forumbuddy/models"
)

type PostRepository interface {
	GetUserById(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(username, passwordHash string) (*models.User, error)
	VerifyUser(username, password string) (*models.User, error)
}
