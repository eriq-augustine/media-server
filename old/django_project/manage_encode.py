'''
from media_server import settings
from django.core.management import setup_environ

setup_environ(settings)
'''

from mediaserver.models import Cache, EncodeQueue

#cache_items = Cache.objects.all()
#
#for item in cache_items:
#    print item.hash
