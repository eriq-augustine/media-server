from django.http import Http404, HttpResponse, HttpResponseRedirect
from django.conf import settings

import os.path

def index(request):
    return HttpResponse("TODO(eriq): index.")

def home(request):
    return HttpResponse("TODO(eriq): home.")

def browse(request, path):
    print path
    print settings.ROOT_DIR
    print settings.CACHE_DIR

    target_path = os.path.join(settings.ROOT_DIR, path)
    target_path = os.path.realpath(target_path)

    if not target_path.startswith(settings.ROOT_DIR):
        raise Http404

    # TODO(eriq): Redirect to viewing this file.
    if not os.path.isdir(target_path):
        raise Http404

    for dir_ent in os.listdir(target_path):
        print dir_ent

    return HttpResponse("TODO(eriq): browse: {}".format(path))
