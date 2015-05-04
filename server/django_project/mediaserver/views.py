from fileutils import EXTENSIONS, Path, UnsafePath, OutOfRootException, PathDoesntExistsException

import cache
import fileutils
import encode

from django.conf import settings
from django.core.servers.basehttp import FileWrapper
from django.core.urlresolvers import reverse
from django.http import Http404, HttpResponse, HttpResponseRedirect, StreamingHttpResponse
from django.shortcuts import render, redirect

from operator import itemgetter
import io
import json
import os
import sys

# TODO(eriq): ?
def index(request):
   return HttpResponseRedirect(reverse('browse'))

# TODO(eriq): Playlists and what-ot.
def home(request):
   return HttpResponseRedirect(reverse('browse'))

def fetch_cache(request, urlpath):
   # HACK(eriq): Don't use unsafe path.
   #  It is safe, just not rooted in the media dir.

   path = UnsafePath.from_urlpath(urlpath)

   ext = path.ext()

   if not ext in EXTENSIONS:
      raise Http404

   mime = EXTENSIONS[ext]['mime']

   if not mime:
      raise Http404

   address = "{}/{}".format(settings.CACHE_SERVE_BASE, urlpath)

   return redirect(address)

def raw(request, urlpath):
   try:
      path = Path.from_urlpath(urlpath)
   except OutOfRootException:
      raise Http404
   except PathDoesntExistsException:
      raise Http404

   if path.is_dir():
      return HttpResponseRedirect(reverse('browse', args = [path.urlpath()]))

   ext = path.ext()

   if not ext in EXTENSIONS:
      raise Http404

   mime = EXTENSIONS[ext]['mime']

   if not mime:
      raise Http404

   address = "{}/{}".format(settings.MEDIA_SERVE_BASE, path.urlpath())

   return redirect(address)

def view(request, urlpath):
   try:
      path = Path.from_urlpath(urlpath)
   except OutOfRootException:
      raise Http404
   except PathDoesntExistsException:
      raise Http404

   if path.is_dir():
      return HttpResponseRedirect(reverse('browse', args = [path.urlpath()]))

   ext = path.ext()

   context = {
      'full_path': path.syspath(),
      'path': path.urlpath(),
      'name': path.display_name(),
      'parent': path.parent().urlpath(),
      'type': ext,
      'breadcrumbs': fileutils.build_breadcrumbs(path),
   }

   # Add a link to direct download the file.
   download_url = reverse('raw', args = [path.urlpath()])
   context['file_context_action'] = {
      'text': 'Direct Download',
      'action': "window.location.href = '" + download_url + "';"
   }

   if ext not in EXTENSIONS:
      return render(request, 'mediaserver/unsupported_file.html', context)

   # Check if this format needs a cache.
   if 'cache' in EXTENSIONS[ext]:
      # Check if the cache is ready.
      cache_status, cache_data = cache.fetch_cache_data(path)

      if cache_status == cache.CACHE_STATUS_READY:
         cache.prep_context(path, context, cache_data)
         return render(request, EXTENSIONS[ext]['template'], context)
      elif cache_status == cache.CACHE_STATUS_ENCODING:
         return render(request, EXTENSIONS[ext]['cache']['place_holder_template'], context)
      elif cache_status == cache.CACHE_STATUS_ERROR:
         # TODO(eriq): Error template?
         raise Exception('Cache Error')
      else:
         raise Exception('Unknown cache status: ' + cache_status)
   else:
      context['mime'] = EXTENSIONS[ext]['mime']
      return render(request, EXTENSIONS[ext]['template'], context)

# TODO(eriq): Need to convert paths back to url paths before render.
#  Windows will have a problem.
def browse(request, urlpath = ''):
   try:
      path = Path.from_urlpath(urlpath)
   except OutOfRootException:
      raise Http404
   except PathDoesntExistsException as err:
      raise Http404

   if path.is_file():
      return HttpResponseRedirect(reverse('view', args = [path.urlpath()]))

   dirs = []
   files = []
   has_image = False

   for dir_ent in os.listdir(path.syspath()):
      dir_ent_path = path.safe_join(dir_ent)

      if dir_ent_path.is_hidden():
        continue

      if dir_ent_path.is_dir():
         dirs.append({'name': dir_ent_path.display_name(),
                      'path': dir_ent_path.urlpath()})
      else:
         ext = dir_ent_path.ext()

         if len(ext) == 0:
            ext = 'txt'

         if ext in EXTENSIONS and EXTENSIONS[ext].get('mime', '').startswith('image/'):
            has_image = True

         files.append({'path': dir_ent_path.urlpath(),
                       'name': dir_ent_path.display_name(),
                       'type': ext})

   context = {
      'full_path': path.syspath(),
      'path': path.urlpath(),
      'dir_name': path.display_name(),
      'parent': path.parent().urlpath(),
      'dirs': sorted(dirs, key = itemgetter('name')),
      'files': sorted(files, key = itemgetter('name')),
      'breadcrumbs': fileutils.build_breadcrumbs(path),
   }

   # Add a gallery button if there are images in this directory.
   if has_image:
      gallery_url = reverse('gallery', args = [path.urlpath()])
      context['file_context_action'] = {
         'text': 'View As Gallery',
         'action': "window.location.href = '" + gallery_url + "';"
      }

   return render(request, 'mediaserver/browse.html', context)

# Make a gallery with the images in the current directory.
def gallery(request, urlpath = ''):
   try:
      path = Path.from_urlpath(urlpath)
   except OutOfRootException:
      raise Http404
   except PathDoesntExistsException as err:
      raise Http404

   if path.is_file():
      return HttpResponseRedirect(reverse('view', args = [path.urlpath()]))

   images = []

   for dir_ent in os.listdir(path.syspath()):
      dir_ent_path = path.safe_join(dir_ent)

      if dir_ent_path.is_hidden():
         continue

      if dir_ent_path.is_dir():
         continue

      ext = dir_ent_path.ext()

      if ext not in EXTENSIONS or not EXTENSIONS[ext].get('mime', '').startswith('image/'):
         continue

      images.append({'path': dir_ent_path.urlpath(),
                     'name': dir_ent_path.display_name(),
                     'type': ext})

   context = {
      'full_path': path.syspath(),
      'path': path.urlpath(),
      'dir_name': path.display_name(),
      'parent': path.parent().urlpath(),
      'images': sorted(images, key = itemgetter('name')),
      'breadcrumbs': fileutils.build_breadcrumbs(path),
   }

   return render(request, 'mediaserver/gallery.html', context)
