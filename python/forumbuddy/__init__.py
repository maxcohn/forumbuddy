from dataclasses import dataclass


import dataclasses
import json
import secrets
import functools

from flask import Flask, request, render_template, session
from flask.helpers import url_for
from flask_session import Session
import argon2
from werkzeug.utils import redirect

from .models import Comment, Post, User


"""
#TODO: Make authentication

    switch session to redis

    /register - GET, POST

    /post - GET, POST
        Make get look better
        Make postable to create new post
    
    /post/:id - comment on post

    /comment/id - GET, POST
        create comment with parent being `id`
        view comment with `id`
"""

app = Flask('asd', template_folder='templates')
app.config['SECRET_KEY'] = secrets.token_bytes(32)

#TODO: add more cookie protection
app.config['SESSION_FILE_DIR'] = '/tmp/flask-session'
app.config['SESSION_COOKIE_NAME'] = 'cs'
app.config['SESSION_TYPE'] = 'filesystem' #TODO: change this to redis

Session(app)

def login_required(f):
    @functools.wraps(f)
    def inner(*args, **kwargs):
        if 'uid' not in session:
            return redirect(url_for('get_login', next=request.url))
        return f(*args, **kwargs)
    return inner

@app.route('/', methods=['GET'])
def index():
    posts = Post.get_recent_posts(15)

    user = None
    if 'uid' in session:
        user = User.get_by_id(session['uid'])

    return render_template('index.html', posts=posts, user=user)

@app.route('/login', methods=['GET'])
def get_login():
    return render_template('login.html')

@app.route('/login', methods=['POST'])
def post_login():
    if 'uid' in session:
        return 'Already logged in'

    username = request.form['username']
    password = request.form['password']

    ph = argon2.PasswordHasher()

    try:
        hash = User.get_password_hash(username)

        # Verify password, raises exception if wrong.
        ph.verify(hash, password)

        # check the hash's parameters and if outdated, rehash the user's password in the database.
        if ph.check_needs_rehash(hash):
            User.set_password_hash(username, ph.hash(password))
    except:
        return 'Invalid username or password', 401

    # Now that we know the user has the correct credentials, we can create their session
    session['uid'] = User.get_by_username(username).uid

    print(session)
    return render_template('login.html')

@app.route('/logout', methods=['GET'])
def logout():
    session.clear()
    redirect(url_for('index'))

@app.route('/user/<int:user_id>', methods=['GET'])
def user_page(user_id: int):
    try:
        return User.get_by_id(user_id).as_json()
    except:
        return '', 404

@app.route('/user', methods=['POST'])
def create_user():
    return 'asd', 203

@app.route('/post/<int:post_id>', methods=['GET'])
@login_required
def post_page(post_id: int):
    #TODO: better error handling
    try:
        post = Post.get_post_and_comments(post_id)
        return render_template('post.html', post=post)
    except:
        return '', 404

@app.route('/comment', methods=['POST'])
def create_comment():
    pass#TODO: from form

#print(json.dumps(dataclasses.asdict(Post.get_post_and_comments(1)), indent=4, default=str))

def no_defaults(pairs):
    ret = {}
    for key, value in pairs:
        if value is not None:
            ret[key] = value

    return ret

def in_json(pairs):
    ret = {}
    dataclasses.fields()
    for key, value in pairs:
        if value is not None:
            ret[key] = value

    return ret

#app.run('localhost', 5000)