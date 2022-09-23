package repos

import (
	"database/sql"
	"errors"
	"forumbuddy/models"

	"github.com/alexedwards/argon2id"
	"github.com/jmoiron/sqlx"
)

type UserRepositorySql struct {
	DB *sqlx.DB
}

func (repo *UserRepositorySql) CreateNewUser(username, passwordHash string) (*models.User, error) {
	newUser := new(models.User)
	err := repo.DB.Get(newUser, `
		WITH user_ins AS (
			INSERT INTO users
				(username)
			VALUES
				($1)
			RETURNING uid, username
		),
		hash_ins AS (
			INSERT INTO user_hashes
				(uid, password_hash)
			VALUES
				((SELECT uid FROM user_ins), $2)
			RETURNING uid
		)
		SELECT uid, username FROM user_ins
	`, username, passwordHash)

	// Check if the user insert failed
	if err != nil {
		//TODO: log error
		return nil, err
	}

	return newUser, nil
}

func (repo *UserRepositorySql) GetUserByUsername(username string) (*models.User, error) {
	user := new(models.User)
	err := repo.DB.Get(user, `
		SELECT uid, username, created_at
		FROM users
		WHERE username = $1
	`, username)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepositorySql) GetUserById(id int) (*models.User, error) {
	user := new(models.User)
	err := repo.DB.Get(user, `
		SELECT uid, username, created_at
		FROM users
		WHERE uid = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

//TODO: shoudl this take teh raw password isntead so we can force it to hash it?
//TODO: separate into a verify and a get instead of doing it in one func
func (repo *UserRepositorySql) VerifyUserPassword(username, password string) (*models.User, error) {
	// Get the hash from the DB for this user
	var passwordHash string
	err := repo.DB.Get(&passwordHash, `
		SELECT uh.password_hash
		FROM users AS u, user_hashes AS uh
		WHERE u.uid = uh.uid
			AND u.username = $1
	`, username)

	if err == sql.ErrNoRows { //TODO: different response for no match?
		return nil, err
	} else if err != nil {
		return nil, err
	}

	// Verify password matches the stored password hash
	match, err := argon2id.ComparePasswordAndHash(password, passwordHash)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, errors.New("password hash didn't match")
	}

	// Now that we know the hashes match, query the user
	return repo.GetUserByUsername(username)
}
