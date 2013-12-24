from gevent import monkey
monkey.patch_all()

from mediaserver.models import EncodeQueue, Cache
from mediaserver.encode import hash_path
from mediaserver.fileutils import Path, UnsafePath
from django.conf import settings

#from geventwebsocket import Resource
#from geventwebsocket.server import WebSocketServer

import json
import random
import threading

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

      while True:
         message = ws.receive()
         if message is None:
            break
         else:
            message = json.loads(message)

            r  = "I have received this message from you : %s" % message
            r += "<br>Glad to be your webserver."
            ws.send(json.dumps({'output': r}))

      del sockets[socket_id]

def run():
   global sockets
   sockets = {}

   server = pywsgi.WSGIServer(("", 6060), WebSocketApp(), handler_class=WebSocketHandler)

   # Start the server on it's own thread.
   server_thread = threading.Thread(target = server.serve_forever)
   # Exit the server thread when the main thread terminates
   server_thread.daemon = True
   server_thread.start()

   while True:
      #TEST
      print len(sockets)

      sleep(5)
