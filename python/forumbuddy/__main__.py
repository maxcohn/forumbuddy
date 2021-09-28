from dataclasses import dataclass
from . import db
from .models import Comment, Post, User

import dataclasses
import json

#print(json.dumps(dataclasses.asdict(Post.get_post_and_comments(1)), indent=4, default=str))

def no_defaults(pairs):
    ret = {}
    for key, value in pairs:
        if value is not None:
            ret[key] = value

    return ret
    

print(json.dumps(dataclasses.asdict(Post.get_post_and_comments(1), dict_factory=no_defaults), indent=4, default=str))