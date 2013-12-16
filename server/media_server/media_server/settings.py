"""
Django settings for media_server project.

For more information on this file, see
https://docs.djangoproject.com/en/1.6/topics/settings/

For the full list of settings and their values, see
https://docs.djangoproject.com/en/1.6/ref/settings/
"""

import multiprocessing

# Build paths inside the project like this: os.path.join(BASE_DIR, ...)
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
ROOT_DIR = os.path.abspath(os.path.realpath('/media/media/bittorent/downloads'))

CACHE_DIR = os.path.abspath(os.path.realpath('cache'))
TEMP_CACHE_DIR = os.path.join(CACHE_DIR, 'temp')
ENCODE_CACHE_DIR = os.path.join(CACHE_DIR, 'encode')

# TODO(eriq): This should be set to the proper deployed url.
# FILE_SERVER = '127.0.0.1:3030'
FILE_SERVER = '192.168.1.169:3030'
CACHE_SERVER = '192.168.1.169:4040'

FFMPEG_PATH = '/usr/bin/ffmpeg'
ENCODING_THREADS = '{}'.format(multiprocessing.cpu_count())
