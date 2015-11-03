from gevent import monkey
monkey.patch_all()

from mediaserver.fileutils import Path, write_pid, mkdir_p
from django.conf import settings

import daemon
import json
import lockfile
import os, os.path
import random
import re
import signal
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
                     limit = 20):
   rtn = []

   try:
      cursor = conn.cursor()
      query = "SELECT hash, src, strftime('%H:%M -- %Y-%m-%d', {}, 'localtime') FROM {} ORDER BY {} {}, {} LIMIT {}"
      cursor.execute(query.format(time_field, table,
                                  time_field, sort_order,
                                  second_sort, limit))

      rows = cursor.fetchall()
      for row in rows:
         path = Path.from_abs_syspath(row[1])
         rtn.append({'hash': row[0],
                     'name': path.display_name(),
                     'path': path.urlpath(),
                     'time': row[2]})
   except Exception as ex:
      print 'Error fetching cache: {}'.format(ex)
      pass

   return rtn

def get_cache(conn):
   return get_general_info(conn, 'mediaserver_cache', 'cache_time', 'DESC', 'urlpath')

def get_queue(conn):
   return get_general_info(conn, 'mediaserver_encodequeue', 'queue_time')

def extract_total_time(info_path):
   info_file = open(info_path, 'r')
   with info_file:
      json_info = json.load(info_file)

   return int(float(json_info['format']['duration']))

# Just let it throw on io error.
def extract_encode_time(progress_path):
   progress_file = open(progress_path, 'r')
   # Expecting 165ish characters every set.
   # 2 == SEEK_END
   progress_file.seek(-200, 2)

   data = str(progress_file.read(200)).replace("\n", ' ')
   match = re.search(r'out_time=(\d\d):(\d\d):(\d\d).(\d+)', data)

   if match != None:
      return int(match.group(1)) * 3600 + int(match.group(2)) * 60 + int(match.group(3))

   return None

# Look for .info and .progress files showing the status of an encode.
def check_progress(info):
   times = {}

   mkdir_p(settings.PROGRESS_CACHE_DIR)

   progresses = {}
   for dir_ent in os.listdir(settings.PROGRESS_CACHE_DIR):
      dir_ent = os.path.abspath(os.path.join(settings.PROGRESS_CACHE_DIR, dir_ent))
      ext = os.path.splitext(dir_ent)[1]
      hash = os.path.splitext(os.path.basename(dir_ent))[0]

      if ext == '.info':
         if not hash in progresses:
           progresses[hash] = {}
         progresses[hash]['info'] = dir_ent
      elif ext == '.progress':
         if not hash in progresses:
           progresses[hash] = {}
         progresses[hash]['progress'] = dir_ent

   for hash in progresses:
      try:
         if 'info' in progresses[hash] and 'progress' in progresses[hash]:
            total_time = extract_total_time(progresses[hash]['info'])
            current_encode_time = extract_encode_time(progresses[hash]['progress'])

            if total_time != None and current_encode_time != None:
               times[hash] = {'total': total_time,
                              'current': current_encode_time}
      except Exception as ex:
        # File could have been removed by encoder.
        pass

   for queued in info['encode_queue']:
      if queued['hash'] in times:
         queued['progress'] = times[queued['hash']]

def get_cache_queue_details(conn):
   info = {'encode_queue': get_queue(conn),
           'recent_cache': get_cache(conn)}

   if len(info['encode_queue']) > 0:
      check_progress(info)

   return info

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

def process_websocket():
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

def run():
   # See notes in manage_encode.py about the DaemonContext.

   if os.path.exists(settings.WEBSOCKET_PID_FILE):
      return

   # daemonize
   context = daemon.DaemonContext(
      working_directory = settings.BASE_DIR,
      stdout = sys.stdout,
      stderr = sys.stderr,
      detach_process = True,
      pidfile = open(settings.WEBSOCKET_PID_FILE, 'w'),
   )

   with context:
      write_pid(settings.WEBSOCKET_PID_FILE)
      process_websocket()
