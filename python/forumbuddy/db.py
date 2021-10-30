'''Wrapper around essential DB access functionality'''

from typing import Any, Tuple, List, Dict
import psycopg2
from psycopg2 import pool, extras

_conn_pool = psycopg2.pool.ThreadedConnectionPool(5, 20, user='postgres', password='password', host='127.0.0.1', port='5432', database='postgres', cursor_factory=psycopg2.extras.DictCursor, connect_timeout=3)

def query(query: str, data: Tuple) -> List[Dict[str, Any]]:
    '''Query the database via a connection pool and return all results

    Arguments:
        query - SQL query to run
        data - Parameters for the query

    Returns:
        List[Dict] - Query results as a list of dictionaries
    '''
    conn = _conn_pool.getconn()
    res = None
    with conn.cursor() as cur:
        cur.execute(query, data)
        res = cur.fetchall()
    _conn_pool.putconn(conn)

    if res is None:
        raise Exception('Query failed') #TODO: make custom exception

    return [dict(row) for row in res]

def query_one(query: str, data: Tuple) -> Dict[str, Any]:
    '''Query the database via a connection pool and return all results

    Arguments:
        query - SQL query to run
        data - Parameters for the query

    Returns:
        List[Dict] - Query results as a list of dictionaries
    '''
    conn = _conn_pool.getconn()
    res = None
    with conn.cursor() as cur:
        cur.execute(query, data)
        res = cur.fetchone()
    _conn_pool.putconn(conn)

    if res is None:
        raise Exception('Query failed') #TODO: make custom exception

    return dict(res)

def execute(query: str, data: Tuple):
    '''Run a query via a connection pool and don't return results
    
    Arguments:
        query - SQL query to run
        data - Parameters for the query
    '''

    conn = _conn_pool.getconn()
    res = None
    with conn.cursor() as cur:
        cur.execute(query, data)
    _conn_pool.putconn(conn)
