from . import db

x = db.query('SELECT * FROM posts WHERE pid = %s', (1,))[0]

print(x)