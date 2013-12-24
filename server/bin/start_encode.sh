#!/bin/bash

PROJECT_ROOT=/home/eriq/media-server/server

DJANGODIR=$PROJECT_ROOT/django_project  # Django project directory
DJANGO_SETTINGS_MODULE=media_server.settings  # which settings file should Django use

ENCODE_PID=$DJANGODIR/encode.pid.lock
ENCODE_ERROR_BASE=$PROJECT_ROOT/logs/encode

cd $PROJECT_ROOT
./bin/removePartialEncodes.sh

# Activate the virtual environment
cd $DJANGODIR
source $PROJECT_ROOT/bin/activate
export DJANGO_SETTINGS_MODULE=$DJANGO_SETTINGS_MODULE
export PYTHONPATH=$DJANGODIR:$PYTHONPATH

# Remove the encode pid file (if it exists).
test -f $ENCODE_PID && rm -f $ENCODE_PID

# Start the encoder.
python manage.py runscript manage_encode > ${ENCODE_ERROR_BASE}.out 2> ${ENCODE_ERROR_BASE}.err
