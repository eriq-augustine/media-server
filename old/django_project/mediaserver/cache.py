# Cache management

from fileutils import EXTENSIONS, Path, UnsafePath, touch, mkdir_p
from encode import hash_path, try_queue

import glob
import os, os.path
import re
import subprocess
import sys

from django.conf import settings

CACHE_STATUS_READY = 0
CACHE_STATUS_ENCODING = 1
CACHE_STATUS_ERROR = 2

# The main interface for the view.
# Returns a tuple containing:
#   1. The CACHE_STATUS
#   2. The cache (None on error)
def fetch_cache_data(path):
   ext = path.ext()
   cache = {}

   encode_queued = False

   if (('contents' not in EXTENSIONS[ext]['cache']) or (not EXTENSIONS[ext]['cache']['contents'])):
      return CACHE_STATUS_READY, cache

   path_hash = hash_path(path)
   cache_path = ensure_cache(path_hash)

   for cache_contents in EXTENSIONS[ext]['cache']['contents']:
      if cache_contents == 'subtitles':
         cache['subtitles'] = fetch_subtitles(path, path_hash, cache_path)
      elif cache_contents == 'poster':
         cache['poster'] = fetch_poster(path, path_hash, cache_path)
      elif cache_contents == 'video_encode':
         encoded_video = fetch_encode(path, path_hash, cache_path)

         if encoded_video is None:
            encode_queued = True
         else:
            cache['encoded_video'] = encoded_video
      else:
         raise Exception('Unknown cache content: "{}"'.format(cache_content))

   if encode_queued:
      return CACHE_STATUS_ENCODING, cache

   return CACHE_STATUS_READY, cache

def ensure_cache(path_hash):
   syspath = os.path.join(settings.CACHE_DIR, path_hash)

   mkdir_p(syspath)

   return UnsafePath.from_abs_syspath(syspath)

def create_cache_url(path_hash, filename):
   return '/cache/{}/{}'.format(path_hash, filename)

def fetch_poster(path, path_hash, cache_path):
   poster_path = cache_path.join(path_hash + '_poster.png')

   if not os.path.exists(poster_path.syspath()):
      generate_poster(path, path_hash, cache_path)

   return create_cache_url(path_hash, path_hash + '_poster.png')

def generate_poster(path, path_hash, cache_path):
   args = [settings.EXTRACT_POSTER_PATH]
   args += [path.syspath()]
   args += [path_hash]
   args += [cache_path.syspath()]

   subprocess.call(args, stdout = sys.stdout, stderr = sys.stderr)

def fetch_subtitles(path, path_hash, cache_path):
   sub_done_path = cache_path.join(path_hash + '_subtitles.done')

   if not os.path.exists(sub_done_path.syspath()):
      generate_subtitles(path, path_hash, cache_path)

   subs = get_subtitles_from_cache(path_hash, cache_path)
   touch(sub_done_path.syspath())

   return subs

def generate_subtitles(path, path_hash, cache_path):
   args = [settings.EXTRACT_SUBTITLES_PATH]
   args += [path.syspath()]
   args += [path_hash]
   args += [cache_path.syspath()]

   subprocess.call(args, stdout = sys.stdout, stderr = sys.stderr)

def get_subtitles_from_cache(path_hash, cache_path):
   subs = []

   subtitle_regex = path_hash + r'_([a-zA-Z]+)_(\d+)\.vtt'
   for subtitle_file in glob.glob(cache_path.join(path_hash + '_*.vtt').syspath()):
      basename = os.path.basename(subtitle_file)
      match = re.search(subtitle_regex, basename)

      lang = match.group(1).lower()
      if lang in settings.LANGUAGE_CODES:
         lang = settings.LANGUAGE_CODES[lang]

      short_lang = lang
      if short_lang in settings.REVERSE_LANGUAGE_CODES:
         short_lang = settings.REVERSE_LANGUAGE_CODES[short_lang]

      display = lang + ' ' + match.group(2)

      subs.append({
         'file': basename,
         'url': create_cache_url(path_hash, basename),
         'display': display,
         'lang': short_lang,
         'display_lang': lang
      })

   return subs

def fetch_encode(path, path_hash, cache_path):
   encode_path = cache_path.join(path_hash + '.mp4')
   encode_done_path = cache_path.join(path_hash + '_encode.done')

   if not os.path.exists(encode_done_path.syspath()):
      try_queue(path)
      return None

   return create_cache_url(path_hash, path_hash + '.mp4')

# Prep the template context given the cache.
def prep_context(path, context, cache):
   ext = path.ext()
   context['cache'] = cache

   context['mime'] = EXTENSIONS[ext]['mime']

   if 'encode' in  EXTENSIONS[ext]['cache']:
      context['encode'] = EXTENSIONS[ext]['cache']['encode']

   if 'encoded_video' in cache:
      context['path'] = cache['encoded_video']
