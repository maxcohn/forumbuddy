CREATE TABLE users (
	uid SERIAL PRIMARY KEY,
	username TEXT,
	password_hash TEXT, --TODO: look into this
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--TODO: FTS index on title and body
--do we want these together, or separate, so we can prioritze titles?

CREATE TABLE posts (
    pid SERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

	--FULL TEXT INDEX (title, body)
);


CREATE TABLE comments (
	cid SERIAL PRIMARY KEY,
	pid INT REFERENCES posts(pid),
	uid INT REFERENCES users(uid),
	body TEXT,
	parent INT REFERENCES comments(cid),

	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

	-- FULL TEXT INDEX body
	-- foreign key for uid, pid, and parent
);


-- procedure for getting comment tree?

CREATE TABLE post_votes (
	vid INT SERIAL PRIMARY KEY,
	uid INT,
	cid INT
);

CREATE TABLE comment_votes (
	vid INT SERIAL PRIMARY KEY,
	uid INT,
	pid INT
);
