#!/bin/sh

PROJECT_ROOT=/home/eriq/media-server/server

DJANGODIR=$PROJECT_ROOT/django_project  # Django project directory
DJANGO_SETTINGS_MODULE=media_server.settings  # which settings file should Django use

WEBSOCKET_PID=$PROJECT_ROOT/run/websocket.pid
WEBSOCKET_LOG_BASE=$PROJECT_ROOT/logs/websocket

# If the pid file exists and is still valid, then just exit.
if [ -f $WEBSOCKET_PID ] ; then
   ps p `cat ${WEBSOCKET_PID}` | grep websocket > /dev/null 2> /dev/null
   if [ $? -eq 0 ] ; then
      exit 0
   fi

   rm -f $WEBSOCKET_PID
fi

# Activate the virtual environment
cd $DJANGODIR
source $PROJECT_ROOT/bin/activate
export DJANGO_SETTINGS_MODULE=$DJANGO_SETTINGS_MODULE
export PYTHONPATH=$DJANGODIR:$PYTHONPATH

# Start the websocket
python manage.py runscript websocket > ${WEBSOCKET_LOG_BASE}.out 2> ${WEBSOCKET_LOG_BASE}.err

sleep 1
