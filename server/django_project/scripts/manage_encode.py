from django.conf import settings
from django.db.models import Sum
from mediaserver.models import EncodeQueue, Cache
from mediaserver.encode import hash_path
from mediaserver.cache import create_cache_url
from mediaserver.fileutils import UnsafePath, write_pid, mkdir_p, touch

from datetime import datetime
import daemon
import glob
import io
import json
import lockfile
import os, os.path
import re
import shutil
import signal
import subprocess
import sys
import time

from multiprocessing import Process

def run():
   # The pidfile option of the DaemonContext does not seem as useful as first glance.
   #  It does not put the pid in the file.
   #  It does not exit if the file exists.
   # So, manage by hand.
   if os.path.exists(settings.ENCODE_PID_FILE):
      return

   # daemonize
   context = daemon.DaemonContext(
      working_directory = settings.BASE_DIR,
      stdout = sys.stdout,
      stderr = sys.stderr,
      detach_process = True,
      # Locking (locking) does not really seem to help.
      # Also, we want systemd to have access to the pid file.
      pidfile = open(settings.ENCODE_PID_FILE, 'w'),
   )

   with context:
      write_pid(settings.ENCODE_PID_FILE)
      process_queue()

def process_queue():
   # HACK(eriq): Just testing, select later.
   while (True):
      while (has_next_encode()):
         manage_cache_size()
         next_encode()
      time.sleep(5)

def has_next_encode():
   return EncodeQueue.objects.count() > 0

def next_encode():
   try:
      encode_task = EncodeQueue.objects.earliest('queue_time')
   except EncodeQueue.DoesNotExist:
      return False

   src_path = UnsafePath.from_abs_syspath(encode_task.src)
   new_path = maybe_encode(src_path)
   encode_task.delete()

   hash = hash_path(src_path)

   encode_path = get_encode_cache_path(src_path, hash)

   new_cache_item = Cache(src = encode_task.src,
                          hash = hash,
                          urlpath = create_cache_url(hash, hash + '.mp4'),
                          bytes = os.path.getsize(encode_path.syspath()))
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
   cache_dir = target_path.parent()

   # Make a directory for the encode.
   mkdir_p(cache_dir.syspath())

   # Copy the file from its source into the cache.
   copy_to_cache(path, temp_path)

   # Encode the file from the cache.
   encode_file(temp_path, original_hash, target_path, cache_dir)

   # Drop the done file.
   touch(cache_dir.join(original_hash + '_encode.done').syspath())

   # rm the temp file.
   os.remove(temp_path.syspath())

def get_encode_cache_path(path, original_hash):
    # <cache dir>/<hash>/<hash>.mp4
   syspath = os.path.join(settings.CACHE_DIR, original_hash, original_hash + '.mp4')
   new_path = UnsafePath.from_abs_syspath(syspath)
   return new_path

def get_temp_cache_path(path, original_hash):
   # The copy does not yet exist and will probably be outside the root,
   #  need an unsafe path.
   syspath = os.path.join(settings.CACHE_DIR, original_hash, original_hash + '.' + path.ext())
   new_path = UnsafePath.from_abs_syspath(syspath)
   return new_path

# Return the new path.
def copy_to_cache(path, temp_path):
   shutil.copyfile(path.syspath(), temp_path.syspath())

def create_info_file(path, info_file):
    args = [settings.FFPROBE_PATH]
    args += ['-v', 'quiet']
    args += ['-print_format', 'json']
    args += ['-show_format']
    args += ['-i', path.syspath()]

    out_file = open(info_file, 'w')
    with out_file:
        subprocess.call(args, stdout = out_file, stderr = sys.stderr)

# Need to keep the original hash, because |path| points to the temp cache
#  not the source file.
def encode_file(path, original_hash, target_path, cache_dir):
   # These files are for the monitoring the status of encodes.
   progress_file = os.path.join(settings.PROGRESS_CACHE_DIR, "{}.progress".format(original_hash))
   info_file = os.path.join(settings.PROGRESS_CACHE_DIR, "{}.info".format(original_hash))

   create_info_file(path, info_file)

   args = [settings.WEBENCODE_PATH]
   args += [path.syspath()]
   args += [original_hash]
   args += [settings.ENCODING_THREADS]
   # Send out the progress.
   args += [progress_file]

   subprocess.call(args, stdout = sys.stdout, stderr = sys.stderr)

    # Remove the info and progress files.
   os.remove(progress_file)
   os.remove(info_file)

def manage_cache_size():
   size_bytes = Cache.objects.aggregate(Sum('bytes'))['bytes__sum']
   max_size_bytes = settings.MAX_CACHE_SIZE_GB * 1024 * 1024 * 1024

   if size_bytes < max_size_bytes:
      return

   lower_size_bytes = settings.CACHE_LOWER_SIZE_GB * 1024 * 1024 * 1024

   cache_objs = Cache.objects.all()
   # [(score, obj), ...]
   cache_scores = [(score_cache_obj(cache_obj), cache_obj) for cache_obj in cache_objs]

   for cache_score, cache_obj in sorted(cache_scores, key = lambda val: val[0], reverse = True):
      if size_bytes < lower_size_bytes:
         break

      size_bytes -= cache_obj.bytes
      remove_from_cache(cache_obj)

# Determine the score for an item in the cache.
# The higher the score, the more likely an item will be removed.
def score_cache_obj(cache_obj):
   last_access = cache_obj.last_access
   now = datetime.utcnow().replace(tzinfo = last_access.tzinfo)
   day_diff = (now - last_access).days

   return ((day_diff + 1) * cache_obj.bytes) / ((cache_obj.hit_count + 1) / 2.0)

def remove_from_cache(cache_obj):
   shutil.rmtree(os.path.join(settings.CACHE_DIR, cache_obj.hash))
   cache_obj.delete()
