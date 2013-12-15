from django.conf import settings

import os
import os.path
import re

class OutOfRootException(Exception):
   pass

class PathDoesntExistsException(Exception):
   path = ''

   def __init__(self, path):
      self.path = path

   def __str__(self):
      return "Path doesn't exist: " + self.path

# Represents a path from the user's perspective.
# Construct with the static methods.
# TODO(eriq): I feel like the comparisons to the ROOT_DIR is a bit fragile.
class Path:
   _abs_syspath = ''
   _is_file = True

   def __init__(self, abs_syspath):
      if not abs_syspath.startswith(settings.ROOT_DIR):
         raise OutOfRootException('Path extends before root: ' + abs_syspath)

      if not os.path.exists(abs_syspath):
         raise PathDoesntExistsException(abs_syspath)

      self._abs_syspath = abs_syspath

      if (os.path.isdir(abs_syspath)):
         self._is_file = False

   @staticmethod
   def from_abs_syspath(syspath):
      return Path(os.path.realpath(syspath))

   @staticmethod
   def from_urlpath(urlpath):
      if urlpath == '':
         urlpath = settings.ROOT_DIR

      # Replace the url slashes with whatever the seperator is on the system.
      syspath = urlpath.replace('/', os.sep)
      syspath = os.path.join(settings.ROOT_DIR, syspath)

      return Path.from_abs_syspath(syspath)

   def syspath(self):
      return self._abs_syspath

   def urlpath(self):
      return self._abs_syspath[(len(settings.ROOT_DIR) + 1):]

   def is_root(self):
      return self._abs_syspath == settings.ROOT_DIR

   # The parent of the root is the root.
   def parent(self):
      if self.is_root():
         parent_path = self._abs_syspath
      else:
         parent_path = os.path.abspath(os.path.join(self._abs_syspath, os.pardir))

      return Path.from_abs_syspath(parent_path)

   def is_dir(self):
      return not self._is_file

   def is_file(self):
      return self._is_file

   def display_name(self):
      if self.is_root():
         return '/'

      return os.path.basename(self._abs_syspath)

   def ext(self):
      if self.is_dir():
         return ''

      return re.sub(r'^\.', '', os.path.splitext(self._abs_syspath)[1])

   def join(self, child):
      return Path.from_abs_syspath(os.path.join(self._abs_syspath, child))

# Start at |path| and go back all the way to root.
def build_breadcrumbs(path):
   crumbs = [{'name': path.display_name(), 'path': path.urlpath()}]

   while not path.is_root():
      path = path.parent()
      crumbs.append({'name': path.display_name(), 'path': path.urlpath()})

   crumbs.reverse()
   return crumbs

EXTENSIONS = {
   # Default none to text.
   '': {'mime': 'text/plain', 'template': 'mediaserver/text_file.html'},
   'txt': {'mime': 'text/plain', 'template': 'mediaserver/text_file.html'},
   'mp3': {'mime': 'audio/mpeg', 'template': 'mediaserver/audio_file.html'},
   'mp4': {'mime': 'video/mp4', 'template': 'mediaserver/video_file.html'},
   'jpg': {'mime': 'image/jpeg', 'template': 'mediaserver/image_file.html'},
   'jpeg': {'mime': 'image/jpeg', 'template': 'mediaserver/image_file.html'},
   'png': {'mime': 'image/png', 'template': 'mediaserver/image_file.html'},
   'gif': {'mime': 'image/gif', 'template': 'mediaserver/image_file.html'},
   'tiff': {'mime': 'image/tiff', 'template': 'mediaserver/image_file.html'},
   'svg': {'mime': 'image/svg+xml', 'template': 'mediaserver/image_file.html'},
}
