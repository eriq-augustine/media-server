from mediaserver.models import EncodeQueue, Cache
from mediaserver.encode import hash_path
from mediaserver.fileutils import Path, UnsafePath
from django.conf import settings

import daemon
import lockfile
import os, os.path
import shutil
import signal
import subprocess
import sys
import time

from multiprocessing import Process

def run():
   # daemonize
   context = daemon.DaemonContext(
      working_directory = settings.BASE_DIR,
      stdout = sys.stdout,
      stderr = sys.stderr,
      detach_process = True,
      pidfile = lockfile.FileLock(settings.ENCODE_PID_FILE),
      signal_map = {
         signal.SIGTTIN: None,
         signal.SIGTTOU: None,
         signal.SIGTSTP: None,
         signal.SIGTERM: None,
      },
   )

   with context:
      process_queue()

def process_queue():
   # HACK(eriq): Just testing, select later.
   while (True):
      while (has_next_encode()):
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
def encode_file(path, original_hash, target_path):
   progress_file = os.path.join(settings.PROGRESS_CACHE_DIR, "{}.progress".format(original_hash))
   info_file = os.path.join(settings.PROGRESS_CACHE_DIR, "{}.info".format(original_hash))

   create_info_file(path, info_file)

   args = [settings.FFMPEG_PATH]
   args += ['-threads', settings.ENCODING_THREADS]
   args += ['-i', path.syspath()]
   args += ['-loglevel', 'warning']
   args += ['-ac', '1']
   args += ['-strict', '-2']
   args += ['-nostats']
   # Send out the progress.
   args += ['-progress', progress_file]

   # Always encode to mp4 since ffmpeg can multi-process well with mp4.
   args += [target_path.syspath()]

   subprocess.call(args, stdout = sys.stdout, stderr = sys.stderr)

    # Remove the info and progress files.
   os.remove(progress_file)
   os.remove(info_file)
