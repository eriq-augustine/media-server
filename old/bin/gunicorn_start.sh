#!/bin/bash

PROJECT_ROOT=/home/eriq/media-server/server

NAME="mediaserver" # Name of the application
DJANGODIR=$PROJECT_ROOT/django_project  # Django project directory
SOCKFILE=$PROJECT_ROOT/run/gunicorn.sock  # we will communicte using this unix socket
USER=media  # the user to run as
GROUP=media  # the group to run as
NUM_WORKERS=1  # how many worker processes should Gunicorn spawn
DJANGO_SETTINGS_MODULE=media_server.settings  # which settings file should Django use
DJANGO_WSGI_MODULE=media_server.wsgi  # WSGI module name

PID_FILE=$PROJECT_ROOT/run/gunicorn.pid

ACCESS_LOG=$PROJECT_ROOT/logs/gunicorn-access.log
ERROR_LOG=$PROJECT_ROOT/logs/gunicorn-error.log

echo "Starting $NAME as `whoami`"

# Activate the virtual environment
cd $DJANGODIR
source $PROJECT_ROOT/bin/activate
export DJANGO_SETTINGS_MODULE=$DJANGO_SETTINGS_MODULE
export PYTHONPATH=$DJANGODIR:$PYTHONPATH

# Create the run directory if it doesn't exist.
RUNDIR=$(dirname $SOCKFILE)
test -d $RUNDIR || mkdir -p $RUNDIR

# Remove the pid file (if it exists).
test -f $PID_FILE && rm -f $PID_FILE

$PROJECT_ROOT/bin/gunicorn ${DJANGO_WSGI_MODULE}:application \
  --name $NAME \
  --workers $NUM_WORKERS \
  --user=$USER --group=$GROUP \
  --debug \
  --log-level debug \
  --access-logfile $ACCESS_LOG \
  --error-logfile $ERROR_LOG \
  --bind=unix:$SOCKFILE \
  --pid $PID_FILE \
  --daemon

sleep 1
