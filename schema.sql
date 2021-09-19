-- User
-- Posts
-- Comments
-- Votes

CREATE TABLE users (
	uid INT SERIAL PRIMARY KEY,
	username TEXT,
	password_hash --TODO: look into this

);

--TODO: FTS index on title and body
--do we want these together, or separate, so we can prioritze titles?

CREATE TABLE posts (
    pid SERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)


CREATE TABLE comments (
	cid INT SERIAL PRIMARY KEY,
	pid INT,
	uid INT,
	body TEXT,
	parent INT,
	FULL TEXT INDEX body
);


-- procedure for getting comment tree?

CREATE TABLE post_votes (
	vid INT SERIAL PRIMARY KEY,
	uid INT,
	cid INT
)

CREATE TABLE comment_votes (
	vid INT SERIAL PRIMARY KEY,
	uid INT,
	pid INT
)
