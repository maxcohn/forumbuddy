from __future__ import annotations
import dataclasses

from datetime import datetime
from typing import Dict, Any, List, Optional, Tuple
from dataclasses import dataclass, field, is_dataclass

import psycopg2

from . import db


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


@dataclass
class User:
    uid: int = 0
    username: str = ''
    password_hash: str = None
    created_at: datetime = None #datetime(1970, 1, 1, 0, 0, 0)
    
    @as_object
    def get_by_id(id: int):
        res = db.query_one('''
            SELECT *
            FROM users
            WHERE uid = %s
            LIMIT 1''', (id,))

        return res

    
@dataclass
class Comment:
    cid: int = 0
    pid: int = 0
    uid: int = 0
    body: str = ''
    parent: Optional[int] = 0
    created_at: datetime = datetime(1970, 1, 1, 0, 0, 0)

    children: List[Comment] = field(default_factory=list)
    user: User = field(default_factory=User)
    
    @as_object
    def get_by_id(id: int):
        res = db.query_one('''
            SELECT *
            FROM comments
            WHERE cid = %s
            LIMIT 1''', (id,))

        return res

    @as_object
    def get_all_from_post(id: int):
        res = db.query('''
            SELECT *
            FROM comments
            WHERE pid = %s
            ''', (1,))

        return res

@dataclass
class Post:

    pid: int = 0
    title: str = ''
    body: str = ''
    created_at: datetime = datetime(1970, 1, 1, 0, 0, 0)

    comments: List[Comment] = field(default_factory=Comment)
    user: User = field(default_factory=User)
    
    @as_object
    def get_by_id(id: int):
        res = db.query_one('''
            SELECT *
            FROM posts
            WHERE pid = %s
            LIMIT 1''', (id,))

        return res

    @classmethod
    def get_post_and_comments(cls, id: int):
        res = db.query_one('''
            SELECT *
            FROM posts
            WHERE pid = %s
            LIMIT 1''', (id,))

        comments = cls.get_comment_tree(id)

        return cls(**res, comments=comments)


    @classmethod
    def get_comment_tree(cls, id: int):
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
            row = dict(row.items())
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