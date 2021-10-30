from __future__ import annotations

import dataclasses
from datetime import datetime
from typing import Dict, Any, List, Optional, Tuple, Union
from dataclasses import dataclass, field, is_dataclass
import json

from . import db


"""
def as_object(func):
    '''Convert the return result of a query into an instance of the retrieved object'''
    def wrapper(cls, *args):
        assert is_dataclass(cls)

        res = func(*args)

        if isinstance(res, psycopg2.extras.DictRow):
            return cls(**res)
        elif type(res) is list:
            return [cls(**r) for r in res]

    return classmethod(wrapper)
"""
def as_object(cls, results: Union[dict, List[Dict[str, any]]]):
    if isinstance(results, dict):
            return cls(**results)
    elif type(results) is list:
        return [cls(**r) for r in results]

class ModelBase:
    #TODO: calculate ahead of time so we don't do all of this every time
    def as_json(self, minified=True):
        exclude_fields = set()
        for field in dataclasses.fields(self):
            if 'in_json' in field.metadata and field.metadata['in_json'] == False:
                exclude_fields.add(field.name)
        
        def factory(pairs):
            ret = {}
            for key, value in pairs:
                if key not in exclude_fields:
                    ret[key] = value

            return ret

        indent = 4
        separators = None
        if minified:
            indent = None
            separators = (',', ':')
        
        return json.dumps(dataclasses.asdict(self, dict_factory=factory), indent=indent, separators=separators, default=str)

@dataclass
class User(ModelBase):
    uid: int = 0
    username: str = ''
    created_at: datetime = datetime(1970, 1, 1, 0, 0, 0)
    
    @classmethod
    def get_by_id(cls, id: int) -> User:
        res = db.query_one('''
            SELECT *
            FROM users
            WHERE uid = %s
            LIMIT 1''', (id,))

        return as_object(cls, res)

    @classmethod
    def get_by_username(cls, username: str) -> User:
        res = db.query_one('''
            SELECT *
            FROM users
            WHERE username = %s
            LIMIT 1''', (username,))

        return as_object(cls, res)

    def get_password_hash(username: str) -> str:
        res = db.query_one('''
            SELECT password_hash
            FROM user_hashes AS uh, users AS u
            WHERE uh.uid = u.uid
                AND u.username = %s''', (username,))

        if res is None:
            raise Exception('No record')
        
        return res['password_hash']

    def set_password_hash(username: str, hash: str):
        db.execute('''
            UPDATE user_hashes AS uh
            SET user_hash = %s
            FROM users AS u
            WHERE uh.uid = u.uid
                AND u.username = %s
        ''', (hash, username))

    
@dataclass
class Comment(ModelBase):
    cid: int = 0
    pid: int = 0
    uid: int = 0
    body: str = ''
    parent: Optional[int] = 0
    created_at: datetime = datetime(1970, 1, 1, 0, 0, 0)

    children: List[Comment] = field(default_factory=list)
    user: User = field(default_factory=User)
    
    @classmethod
    def get_by_id(cls, id: int) -> Comment:
        res = db.query_one('''
            SELECT *
            FROM comments
            WHERE cid = %s
            LIMIT 1''', (id,))

        return as_object(cls, res)

    @classmethod
    def get_all_from_post(cls, id: int) -> List[Comment]:
        res = db.query('''
            SELECT *
            FROM comments
            WHERE pid = %s
            ''', (id,))

        return as_object(cls, res)

@dataclass
class Post(ModelBase):

    pid: int = 0
    title: str = ''
    body: str = ''
    created_at: datetime = datetime(1970, 1, 1, 0, 0, 0)

    comments: List[Comment] = field(default_factory=list)
    user: User = field(default_factory=User)
    
    @classmethod
    def get_by_id(cls, id: int) -> Post:
        res = db.query_one('''
            SELECT *
            FROM posts
            WHERE pid = %s
            LIMIT 1''', (id,))

        return as_object(cls, res)

    @classmethod
    def get_recent_posts(cls, num_posts: int) -> List[Post]:
        res = db.query('''
            SELECT p.pid, p.title, p.body, p.created_at, u.uid, u.username
            FROM posts AS p, users AS u
            WHERE p.uid = u.uid
            ORDER BY created_at DESC
            LIMIT %s''', (num_posts,))

        posts = []
        for row in res:
            user = User(uid=row['uid'], username=row['username'])

            del row['uid']
            del row['username']

            posts.append(cls(**row, user=user))

        return posts

    @classmethod
    def get_post_and_comments(cls, id: int) -> Post:
        res = db.query_one('''
            SELECT p.pid, p.title, p.body, p.created_at, u.username, u.uid
            FROM posts AS p, users AS u
            WHERE p.uid = u.uid
                AND pid = %s
            LIMIT 1''', (id,))

        user = User(uid=res['uid'], username=res['username'])
        comments = cls.get_comment_tree(id)

        del res['uid']
        del res['username']

        return cls(**res, user=user, comments=comments)


    @classmethod
    def get_comment_tree(cls, id: int) -> List[Comment]:
        res = db.query('''
            SELECT c.cid, c.body, c.parent, c.pid, u.uid, u.username
            FROM comments AS c, users AS u
            WHERE c.uid = u.uid
                AND pid = %s
            ORDER BY
                CASE WHEN parent IS NULL THEN 0
                ELSE parent
            END ASC''', (id,))
        
        comment_map: Dict[int, Comment] = {}
        root = []
        for row in res:
            user = User(uid=row['uid'], username=row['username'])
            #del row['uid']
            del row['username']
            cur = Comment(**row, user=user)
            comment_map[cur.cid] = cur

            if cur.parent is None:
                root.append(cur)
            else:
                comment_map[cur.parent].children.append(cur) 

        return root