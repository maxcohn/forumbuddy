--TODO: set auto increment for all of these columns


INSERT INTO users
    (uid, username)
VALUES
    (1, 'foo-bar'),
    (2, 'pg'),
    (3, 'ken'),
    (4, 'dmr');

INSERT INTO user_hashes
    (uid, password_hash)
VALUES
    (1, '$argon2id$v=19$m=65536,t=1,p=2$2mhjjCYDAdJLL8PFC+Gs5w$ZVyL0AJfXi5ps4XbkQ3DudHP/onN0jUMAsvdm/JSB+U'), -- Password = 'password'
    (2, '$argon2id$v=19$m=65536,t=1,p=2$J/9mQ4q9/D3bQd9HuhOBYQ$Bd4lU8EGIRrg7ChZUQSbJutBBIg3ec4Gy+eZXF6x9fQ'), -- Password = 'abc
    (3, '$argon2id$v=19$m=65536,t=1,p=2$HZbWDZ+ORK/N6WjgAy15yQ$TSEtHmsv0CuCqfGW/ubTsSPjXiX+Igwr99LqPRdgDN0'), -- Password = 'greentea'
    (4, '$argon2id$v=19$m=65536,t=1,p=2$vePuEHk/1uVkMMNdl8JVBw$fs5AET22Xd7KL7fQ8PD08mVmpyLFzj/jcO+t+XxRVNc'); -- Password = '!Wh#olaksnd3nkajsnd88?'


INSERT INTO posts
    (pid, uid, title, body)
VALUES
    (1, 2, 'What *is* eating Gilber Grape?', 'We''re honestly not sure, we should probably watch the movie'),
    (2, 3, 'Green is the color of life', 'What about fully ripe bell peppers?');

INSERT INTO comments
    (cid, pid, uid, body, parent)
VALUES
    (1, 1, 1, 'Oh man, you really should watch that movie, Leo does a great job.', NULL),
    (2, 1, 2, 'Really? I wasn''t a huge fan honestly.', 1),
    (3, 1, 3, 'Come on man. What about Johnny Depp in that movie? Killer performance.', 2),
    (4, 1, 4, 'Eh, mediocre.', NULL),
    (5, 2, 1, 'Is there an acutal difference in taste or nutrition depending on the level of ripeness of a bell pepper?', NULL),
    (6, 2, 2, 'I don''t think so', 5),
    (7, 2, 4, 'Never liked bell peppers personally.', NULL);


-- Update all sequence since we added in IDs manually
select setval('users_uid_seq', (select max(uid) from users));
select setval('comments_cid_seq', (select max(cid) from comments));
select setval('posts_pid_seq', (select max(pid) from posts));