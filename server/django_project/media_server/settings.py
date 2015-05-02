"""
Django settings for media_server project.

For more information on this file, see
https://docs.djangoproject.com/en/1.6/topics/settings/

For the full list of settings and their values, see
https://docs.djangoproject.com/en/1.6/ref/settings/
"""

import multiprocessing

import os
BASE_DIR = os.path.dirname(os.path.dirname(__file__))


# Quick-start development settings - unsuitable for production
# See https://docs.djangoproject.com/en/1.6/howto/deployment/checklist/

# SECURITY WARNING: keep the secret key used in production secret!
SECRET_KEY = 'w=@p#=q9ir*31=1rgl-0!i7nb1=^h(s(*@)$@1mjh0vh@h6@30'

# SECURITY WARNING: don't run with debug turned on in production!
DEBUG = True

TEMPLATE_DEBUG = True

ALLOWED_HOSTS = []


# Application definition

INSTALLED_APPS = (
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
    'mediaserver',
    'django_extensions',
)

MIDDLEWARE_CLASSES = (
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
)

ROOT_URLCONF = 'media_server.urls'

WSGI_APPLICATION = 'media_server.wsgi.application'


# Database
# https://docs.djangoproject.com/en/1.6/ref/settings/#databases

SQLITE_DB_PATH = os.path.join(BASE_DIR, 'db.sqlite3')
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.sqlite3',
        'NAME': os.path.join(BASE_DIR, 'db.sqlite3'),
    }
}

#DATABASES = {
#    'default': {
#        'ENGINE': 'django.db.backends.mysql',
#        'NAME': 'FANIME_PANELS',
#        'USER': 'panels',
#        'PASSWORD': 'ILovePanels',
#        'HOST': '',                      # Empty for localhost through domain sockets or '127.0.0.1' for localhost through TCP.
#        'PORT': '',                      # Set to empty string for default.
#    }
#}

# Internationalization
# https://docs.djangoproject.com/en/1.6/topics/i18n/

LANGUAGE_CODE = 'en-us'

TIME_ZONE = 'America/Los_Angeles'

USE_I18N = True

USE_L10N = True

USE_TZ = True

TEMPLATE_DIRS = (
    os.path.join(BASE_DIR, 'templates')
)

# List of callables that know how to import templates from various sources.
TEMPLATE_LOADERS = (
    'django.template.loaders.filesystem.Loader',
    'django.template.loaders.app_directories.Loader',
#     'django.template.loaders.eggs.Loader',
)

# Static files (CSS, JavaScript, Images)
# https://docs.djangoproject.com/en/1.6/howto/static-files/

STATIC_URL = '/static/'

# List of finder classes that know how to find static files in
# various locations.
STATICFILES_FINDERS = (
    'django.contrib.staticfiles.finders.FileSystemFinder',
    'django.contrib.staticfiles.finders.AppDirectoriesFinder',
#    'django.contrib.staticfiles.finders.DefaultStorageFinder',
)

# Project constants.
#ROOT_DIR = os.path.abspath(os.path.realpath('/media/media/bittorent/downloads'))
#ROOT_DIR = os.path.abspath(os.path.realpath('/media/nas'))
# {git root}/server/media
# Need to get the realpath to negotiate symlinks.
#  Otherwise it will be difficult to tell if the request is outside of root.
ROOT_DIR = os.path.realpath(os.path.join(BASE_DIR, os.pardir, 'media'))

BIN_DIR = os.path.realpath(os.path.join(BASE_DIR, os.pardir, 'bin'))

CACHE_DIR = os.path.abspath(os.path.join(BASE_DIR, os.pardir, 'cache'))

PROGRESS_CACHE_DIR = os.path.join(CACHE_DIR, 'progress')

# When cache is larger than this, remove items.
MAX_CACHE_SIZE_GB = 40
# While cache is above this, remove items.
CACHE_LOWER_SIZE_GB = 30

# These are the path to alternate locations served through the webserver.
MEDIA_SERVE_BASE = '/media'
CACHE_SERVE_BASE = '/cache'

FFMPEG_PATH = '/usr/bin/ffmpeg'
FFPROBE_PATH = '/usr/bin/ffprobe'

ENCODE_UTILS_PATH = os.path.join(BIN_DIR, 'encode')
EXTRACT_POSTER_PATH = os.path.join(ENCODE_UTILS_PATH, 'extract_poster')
EXTRACT_SUBTITLES_PATH = os.path.join(ENCODE_UTILS_PATH, 'extract_subtitles')
WEBENCODE_PATH = os.path.join(ENCODE_UTILS_PATH, 'webencode')

ENCODING_THREADS = '{}'.format(multiprocessing.cpu_count())
ENCODE_PID_FILE = os.path.join(BASE_DIR, os.pardir, 'run', 'encode.pid')

WEBSOCKET_PORT = 6060
WEBSOCKET_ADDRESS = "71.84.26.224:{}".format(WEBSOCKET_PORT)
WEBSOCKET_PID_FILE = os.path.join(BASE_DIR, os.pardir, 'run', 'websocket.pid')

# TODO(eriq): Get a full list of codes.
LANGUAGE_CODES = {
   'eng': 'English',
   'en': 'English',
   'english': 'English',
   'jp': 'Japanese',
   'jpn': 'Japanese',
   'fr': 'French',
   'es': 'Spanish',
   'spa': 'Spanish'
}

REVERSE_LANGUAGE_CODES = {
   'English': 'en',
   'Japanese': 'jp',
   'French': 'fr',
   'Spanish': 'es'
}
