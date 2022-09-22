package repos

import (
	"database/sql"
	"fmt"
	"forumbuddy/models"

	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	GetUserById(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(username, passwordHash string) (*models.User, error)
	VerifyUser(username, password string) (*models.User, error)
}

type CommentRepositorySql struct {
	DB *sqlx.DB
}

func (repo *CommentRepositorySql) GetCommentById(id int) (*models.Comment, error) {
	comment := new(models.Comment)

	err := repo.DB.Get(comment, `
		SELECT
			c.cid,
			c.pid,
			c.body,
			c.parent,
			c.created_at,
			u.uid AS "user.uid",
			u.username AS "user.username"
		FROM comments AS c, users AS u
		WHERE c.uid = u.uid
			AND c.cid = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (repo *CommentRepositorySql) CreateNewComment(uid, pid int, parent sql.NullInt64, body string) (int, error) {
	// Insert the comemnt into the DB and get that new comment's ID
	var newCommentId int
	err := repo.DB.Get(&newCommentId, `
		INSERT INTO comments
			(uid, pid, parent, body)
		VALUES
			($1, $2, $3, $4)
		RETURNING cid
	`, uid, pid, parent, body)

	if err != nil {
		fmt.Println(err.Error())
		//TODO: log this?
		return 0, err
	}

	return newCommentId, nil
}
