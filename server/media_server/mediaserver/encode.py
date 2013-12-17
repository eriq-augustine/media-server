# Base for encoding and serving files.

from mediaserver.models import EncodeQueue, Cache
from fileutils import Path, UnsafePath
from django.conf import settings

import hashlib
import os.path
import shutil
import subprocess
import threading

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
      return Cache.objects.get(hash = hash_path(path))
   except Cache.DoesNotExist:
      return None

def queue(path):
   queue_item = EncodeQueue(src = path.syspath(), hash = hash_path(path))
   queue_item.save()

   global encode_thread

   try:
      encode_thread
   except NameError:
      encode_thread = None

   if encode_thread == None:
      encode_thread = threading.Thread(target = process_queue)
      encode_thread.setDaemon(True)
      encode_thread.start()

def process_queue():
   while (has_next_encode()):
      next_encode()

   # HACK(eriq): Race condition here.
   global encode_thread
   encode_thread = None

def hash_path(path):
   md5 = hashlib.md5()
   md5.update(path.syspath())
   return md5.hexdigest()

def has_next_encode():
   return EncodeQueue.objects.count() > 0

def next_encode():
   try:
      encode_task = EncodeQueue.objects.latest('queue_time')
   except EncodeQueue.DoesNotExist:
      return False

   src_path = UnsafePath.from_abs_syspath(encode_task.src)
   new_path = maybe_encode(src_path)
   encode_task.delete()

   # HACK(eriq): The cache urlpath is pretty janky.
   hash = hash_path(src_path)
   new_cache_item = Cache(src = encode_task.src,
                          hash = hash,
                          urlpath = hash + '.mp4')
   new_cache_item.save()

# Maybe encode the file and return the path to the converted file.
def maybe_encode(path):
   original_hash = hash_path(path)

   # Check the cache first
   temp_path = get_temp_cache_path(path, original_hash)
   encode_path = get_encode_cache_path(path, original_hash)

   if not os.path.exists(encode_path.syspath()):
      encode(path, original_hash, temp_path, encode_path)

   return encode_path

def encode(path, original_hash, temp_path, target_path):
   # Copy the file from its source into the cache.
   copy_to_cache(path, temp_path)

   # Encode the file from the cache.
   encode_file(temp_path, original_hash, target_path)

   # rm the temp file.
   os.remove(temp_path.syspath())

def get_encode_cache_path(path, original_hash):
   syspath = os.path.join(settings.ENCODE_CACHE_DIR, original_hash + '.mp4')
   new_path = UnsafePath.from_abs_syspath(syspath)
   return new_path

def get_temp_cache_path(path, original_hash):
   # The copy does not yet exist and will probably be outside the root,
   #  need an unsafe path.
   syspath = os.path.join(settings.TEMP_CACHE_DIR, original_hash + '.' + path.ext())
   new_path = UnsafePath.from_abs_syspath(syspath)
   return new_path

# Return the new path.
def copy_to_cache(path, temp_path):
   shutil.copyfile(path.syspath(), temp_path.syspath())

# Need to keep the original hash, because |path| points to the temp cache
#  not the source file.
def encode_file(path, original_hash, target_path):
   args = [settings.FFMPEG_PATH]
   args += ['-threads', settings.ENCODING_THREADS]
   args += ['-i', path.syspath()]
   args += ['-ac', '1']
   args += ['-strict', '-2']

   # Always encode to mp4 since ffmpeg can multi-process well with mp4.
   args += [target_path.syspath()]

   return subprocess.call(args)
