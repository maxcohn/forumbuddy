package repos

import (
	"forumbuddy/models"

	"github.com/jmoiron/sqlx"
)

type PostRepositorySql struct {
	DB *sqlx.DB
}

//TODO: get comment and children by id

func (repo *PostRepositorySql) GetPostAndCommentsById(id int) (*models.Post, error) {
	post := new(models.Post)
	var comments []models.Comment

	// Query the current post
	err := repo.DB.Get(post, `
		SELECT
			p.pid,
			p.title,
			p.body,
			p.created_at,
			u.uid AS "user.uid",
			u.username AS "user.username"
		FROM posts AS p, users AS u
		WHERE p.uid = u.uid
			AND p.pid = $1
		
	`, id)

	if err != nil {
		return nil, err
	}

	// Query all comments on the post
	err = repo.DB.Select(&comments, `
		SELECT
			c.cid,
			c.body,
			c.parent,
			c.pid,
			u.uid AS "user.uid",
			u.username AS "user.username"
		FROM comments AS c, users AS u
		WHERE c.uid = u.uid
			AND c.pid = $1
		ORDER BY
			CASE WHEN parent IS NULL THEN 0
			ELSE parent
		END ASC
	`, id)

	if err != nil {
		return nil, err
	}

	// Create an empty slice of the comments at the root of the tree. These are pointers since we're going to be updating the slices as we go
	rootComments := make([]*models.Comment, 0)

	// Create an empty mapping from cids to comments. These are points because we're going to be modifying them in our loop
	commentMap := make(map[int]*models.Comment)

	for i, comment := range comments {
		var curComment = &comments[i]

		curComment.Children = make([]*models.Comment, 0)
		commentMap[comment.Cid] = curComment

		if !comment.Parent.Valid {
			// If there is not parent comment, this is at the root
			rootComments = append(rootComments, curComment)
		} else {
			// If there is a parent comment
			parent := commentMap[int(comment.Parent.Int64)]
			parent.Children = append(parent.Children, curComment)
		}
	}

	// Convert all root comments to their values
	for _, c := range rootComments {
		post.Comments = append(post.Comments, *c)
	}

	return post, nil
}

func (repo *PostRepositorySql) CreateNewPost(uid int, title, body string) (int, error) {
	// Insert the post into the DB and get that new post's ID
	var newPostId int
	err := repo.DB.Get(&newPostId, `
		INSERT INTO posts
			(uid, title, body)
		VALUES
			($1, $2, $3)
		RETURNING pid
	`, uid, title, body)

	if err != nil {
		return 0, err
	}

	return newPostId, nil
}

func (repo *PostRepositorySql) GetRecentPosts(limit int) ([]models.Post, error) {
	var posts []models.Post
	err := repo.DB.Select(&posts, `
		SELECT
			p.pid,
			p.title,
			p.created_at,
			u.username AS "user.username"
		FROM posts as p, users AS u
		WHERE p.uid = u.uid
		ORDER by p.created_at desc 
		LIMIT $1
	`, limit)

	if err != nil {
		return nil, err
	}

	return posts, nil
}
