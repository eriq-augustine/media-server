from fileutils import EXTENSIONS, Path, UnsafePath, OutOfRootException, PathDoesntExistsException
import fileutils
import encode

from django.conf import settings
from django.core.servers.basehttp import FileWrapper
from django.core.urlresolvers import reverse
from django.http import Http404, HttpResponse, HttpResponseRedirect, StreamingHttpResponse
from django.shortcuts import render, redirect

from operator import itemgetter
import os
import sys

def index(request):
   return HttpResponse("TODO(eriq): index.")

def home(request):
   return HttpResponse("TODO(eriq): home.")

# TODO(eriq): Error handle this.
def fetch_encode(request, urlpath):
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
      'breadcrumbs': fileutils.build_breadcrumbs(path)
   }

   if ext in EXTENSIONS:
      # Check if this format needs an encode.
      if 'encode' in EXTENSIONS[ext]:
         # Check if the file is cached.
         cache = encode.get_cache(path)
         if cache:
            context['mime'] = EXTENSIONS[ext]['mime']
            context['path'] = cache.urlpath
            context['cache'] = True
            context['encode'] = EXTENSIONS[ext]['encode']
            return render(request, EXTENSIONS[ext]['encode_template'], context)
         elif not encode.is_queued(path):
            encode.queue(path)
            return render(request, EXTENSIONS[ext]['template'], context)
         else:
            return render(request, EXTENSIONS[ext]['template'], context)
      else:
         context['mime'] = EXTENSIONS[ext]['mime']
         return render(request, EXTENSIONS[ext]['template'], context)

   return render(request, 'mediaserver/unsupported_file.html', context)

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
      'breadcrumbs': fileutils.build_breadcrumbs(path)
   }

   return render(request, 'mediaserver/browse.html', context)
