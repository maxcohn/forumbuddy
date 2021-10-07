from dataclasses import dataclass
from . import db
from .models import Comment, Post, User

import dataclasses
import json

from flask import Flask, request, render_template


app = Flask('asd', template_folder='templates')

@app.route('/', methods=['GET'])
def index():
    posts = Post.get_recent_posts(15)

    return render_template('index.html', posts=posts)

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

app.run()