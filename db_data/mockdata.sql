INSERT INTO posts
    (pid, title, body)
VALUES
    (1, 'What *is* eating Gilber Grape?', 'We''re honestly not sure, we should probably watch the movie'),
    (2, 'Green is the color of life', 'What about fully ripe bell peppers?')


INSERT INTO users
    (uid, username, password_hash)
VALUES
    (1, 'foo-bar', ''),
    (2, 'pg', ''),
    (3, 'ken', ''),
    (4, 'dmr', '')

INSERT INTO comments
    (cid, pid, uid, body, parent)
VALUES
    (1, 1, 1, 'Oh man, you really should watch that movie, Leo does a great job.', NULL),
    (2, 1, 2, 'Really? I wasn''t a huge fan honestly.', 1),
    (3, 1, 3, 'Come on man. What about Johnny Depp in that movie? Killer performance.', 2),
    (4, 1, 4, 'Eh, mediocre.', NULL),
    (5, 2, 1, 'Is there an acutal difference in taste or nutrition depending on the level of ripeness of a bell pepper?', NULL),
    (6, 2, 2, 'I don''t think so', 5),
    (7, 2, 4, 'Never liked bell peppers personally.', NULL)