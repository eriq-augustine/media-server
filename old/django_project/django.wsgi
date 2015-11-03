import os
import sys

# TODO(eriq): Not use abs path.
sys.path.append('/home/eriq/media-server/server/media_server')
sys.path.append('/home/eriq/media-server/server/media_server/media_server')
sys.path.append('/home/eriq/media-server/server/media_server/mediaserver')
#sys.path.append(os.path.realpath(os.path.dirname(__file__)))

os.environ['DJANGO_SETTINGS_MODULE'] = 'media_server.settings'

import django.core.handlers.wsgi
application = django.core.handlers.wsgi.WSGIHandler()
