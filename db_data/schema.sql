CREATE TABLE users (
	uid SERIAL PRIMARY KEY,
	username TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--TODO: FTS index on title and body
--do we want these together, or separate, so we can prioritze titles?

CREATE TABLE posts (
    pid SERIAL NOT NULL PRIMARY KEY,
	uid INT NOT NULL REFERENCES users(uid),
    title TEXT NOT NULL,
    body TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

	--FULL TEXT INDEX (title, body)
);


CREATE TABLE comments (
	cid SERIAL PRIMARY KEY,
	pid INT NOT NULL REFERENCES posts(pid),
	uid INT NOT NULL REFERENCES users(uid),
	body TEXT NOT NULL,
	parent INT REFERENCES comments(cid),

	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

	-- FULL TEXT INDEX body
	-- foreign key for uid, pid, and parent
);

CREATE TABLE user_hashes (
	uid INT PRIMARY KEY NOT NULL REFERENCES users(uid),
	password_hash TEXT NOT NULL --TODO: look into this

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
