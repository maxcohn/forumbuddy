package repos

import (
	"database/sql"
	"forumbuddy/models"
)

type CommentRepository interface {
	GetCommentById(id int) (*models.Comment, error)
	CreateNewComment(uid, pid int, parent sql.NullInt64, body string) (int, error)
}
