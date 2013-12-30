# Base for encoding and serving files.

from mediaserver.models import EncodeQueue, Cache
from fileutils import Path, UnsafePath
from django.conf import settings

import hashlib
import os, os.path
import sys

from multiprocessing import Process

# Get some details on the cache and queue for display to the user.
def get_display_details():
   return {'encode_queue': get_encode_queue_items(10),
           'recent_cache': get_recent_cache_items(10)}

def get_encode_queue_items(num):
   rtn = []
   items = EncodeQueue.objects.filter().order_by('queue_time', 'src')[:num]

   for item in items:
      path = Path.from_abs_syspath(item.src)
      rtn.append({'name': path.display_name(),
                  'path': path.urlpath(),
                  'time': item.queue_time})

   return rtn

def get_recent_cache_items(num):
   rtn = []
   items = Cache.objects.filter().order_by('-cache_time', 'urlpath')[:num]

   for item in items:
      path = Path.from_abs_syspath(item.src)
      rtn.append({'name': path.display_name(),
                  'path': path.urlpath(),
                  'time': item.cache_time})

   return rtn

def is_queued(path):
   try:
      encode_queue = EncodeQueue.objects.get(hash = hash_path(path))
      return True
   except EncodeQueue.DoesNotExist:
      return False

def is_cached(path):
   return not get_cache(path) == None

def get_cache(path):
   try:
      cache_obj = Cache.objects.get(hash = hash_path(path))
      cache_obj.hit_count += 1
      cache_obj.save()
      return cache_obj
   except Cache.DoesNotExist:
      return None

def hash_path(path):
   md5 = hashlib.md5()
   md5.update(path.syspath())
   return md5.hexdigest()

def queue(path):
   queue_item = EncodeQueue(src = path.syspath(), hash = hash_path(path))
   queue_item.save()
