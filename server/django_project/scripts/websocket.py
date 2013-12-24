from gevent import monkey
monkey.patch_all()

from mediaserver.fileutils import Path, UnsafePath
from django.conf import settings

import json
import random
import sys
import threading

import sqlite3 as sqlite

from gevent import pywsgi, sleep
from geventwebsocket.handler import WebSocketHandler

class WebSocketApp(object):
   # Establish a websocket.
   def __call__(self, environ, start_response):
      global sockets

      ws = environ['wsgi.websocket']

      # TODO(eriq): Rece condition.
      #  (Not really because of GIL, but we should still treat it as such).
      socket_id = len(sockets)
      sockets[socket_id] = ws

      # Send the initial info.
      send_last_cache_queue_info(ws)

      while True:
         message = ws.receive()
         if message is None:
            break
         else:
            json_message = None
            try:
                json_message = json.loads(message)
            except:
                pass

            # Not really expecting any messages, just ignore any content

      del sockets[socket_id]

# Having a big problem running django db wrappers and gevent threads
#  (which we need to the websocket server).
# So we will have to manage db by hand.
def get_sqlite_conn():
   conn = None

   try:
      conn = sqlite.connect(settings.SQLITE_DB_PATH)
   except sqlite.Error as ex:
      print "SQLite Error {}:".format(ex)

   return conn

def get_general_info(conn, table, time_field,
                     sort_order = '', second_sort = 'src',
                     limit = 10):
   rtn = []

   try:
      cursor = conn.cursor()
      cursor.execute('SELECT src, {} FROM {} ORDER BY {} {}, {} LIMIT {}'.format(time_field,
                                                                                 table,
                                                                                 time_field,
                                                                                 sort_order,
                                                                                 second_sort,
                                                                                 limit))

      rows = cursor.fetchall()
      for row in rows:
         path = Path.from_abs_syspath(row[0])
         rtn.append({'name': path.display_name(),
                     'path': path.urlpath(),
                     'time': row[1]})
   except Exception as ex:
      print 'Error fetching cache: {}'.format(ex)
      pass

   return rtn

def get_cache(conn):
   return get_general_info(conn, 'mediaserver_cache', 'cache_time', 'DESC', 'urlpath')

def get_queue(conn):
   return get_general_info(conn, 'mediaserver_encodequeue', 'queue_time')

def get_cache_queue_details(conn):
      return {'encode_queue': get_queue(conn),
              'recent_cache': get_cache(conn)}

def send_cache_queue_info(info, ws):
    ws.send(json.dumps({'type': 'ENCODE_UPDATE',
                        'info': info}))

# Don't get new info, just send it out.
def send_last_cache_queue_info(ws):
   global last_info

   if not last_info == None:
      send_cache_queue_info(last_info, ws)

# Race condition.
def update_cache_queue_info(conn, ws):
   global last_info
   last_info = get_cache_queue_details(conn)

   # TODO(eriq): Only send if delta.
   send_cache_queue_info(last_info, ws)

def run():
   global sockets
   sockets = {}

   # Only do db access from this thread.
   # So, all clients will be synced on the same data instead of access for each one.
   conn = get_sqlite_conn()

   # Prime the info.
   global last_info
   last_info = get_cache_queue_details(conn)

   if conn == None:
      print 'ERROR: Cannot establish sqlite conn.'
      return

   server = pywsgi.WSGIServer(("", settings.WEBSOCKET_PORT),
                              WebSocketApp(),
                              handler_class = WebSocketHandler)

   # Start the server on it's own thread.
   server_thread = threading.Thread(target = server.serve_forever)
   # Exit the server thread when the main thread terminates
   server_thread.daemon = True
   server_thread.start()

   while True:
      sleep(2)

      for socket_id in sockets:
        # sockets[socket_id].send(json.dumps({'ping': len(sockets)}))
        update_cache_queue_info(conn, sockets[socket_id])

   # Clean it up!
   conn.close()
