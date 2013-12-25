#!/bin/sh

PROJECT_ROOT=/home/eriq/media-server/server

DJANGODIR=$PROJECT_ROOT/django_project  # Django project directory
DJANGO_SETTINGS_MODULE=media_server.settings  # which settings file should Django use

WEBSOCKET_PID=$DJANGODIR/websocket.pid.lock
WEBSOCKET_LOG_BASE=$PROJECT_ROOT/logs/websocket

# Activate the virtual environment
cd $DJANGODIR
source $PROJECT_ROOT/bin/activate
export DJANGO_SETTINGS_MODULE=$DJANGO_SETTINGS_MODULE
export PYTHONPATH=$DJANGODIR:$PYTHONPATH

# Remove the websocket pid file (if it exists).
test -f $WEBSOCKET_PID && rm -f $WEBSOCKET_PID

# Start the websocket
python manage.py runscript websocket > ${WEBSOCKET_LOG_BASE}.out 2> ${WEBSOCKET_LOG_BASE}.err
# gunicorn -k "geventwebsocket.gunicorn.workers.GeventWebSocketWorker" wsgi:websocket_app
