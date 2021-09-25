from dataclasses import dataclass
from . import db
from .models import Comment, Post, User

import dataclasses
import json

#print(json.dumps(dataclasses.asdict(Post.get_post_and_comments(1)), indent=4, default=str))