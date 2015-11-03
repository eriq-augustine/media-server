#!/bin/bash

PROJECT_ROOT=/home/eriq/media-server/server

DJANGODIR=$PROJECT_ROOT/django_project  # Django project directory
DJANGO_SETTINGS_MODULE=media_server.settings  # which settings file should Django use

ENCODE_PID=$PROJECT_ROOT/run/encode.pid
ENCODE_LOG_BASE=$PROJECT_ROOT/logs/encode

# If the pid file exists and is still valid, then just exit.
if [ -f $ENCODE_PID ] ; then
   ps p `cat ${ENCODE_PID}` | grep manage_encode > /dev/null 2> /dev/null
   if [ $? -eq 0 ] ; then
      exit 0
   fi

   rm -f $ENCODE_PID
fi

cd $PROJECT_ROOT
./bin/removePartialEncodes.sh

# Remove any info/progress files
rm -f $PROJECT_ROOT/cache/progress/*

# Activate the virtual environment
cd $DJANGODIR
source $PROJECT_ROOT/bin/activate
export DJANGO_SETTINGS_MODULE=$DJANGO_SETTINGS_MODULE
export PYTHONPATH=$DJANGODIR:$PYTHONPATH

# Start the encoder.
python manage.py runscript manage_encode > ${ENCODE_LOG_BASE}.out 2> ${ENCODE_LOG_BASE}.err

# There is a small race between creating the PID file and systemd checking for it.
#  This little advantage for python should be enough to always win.
sleep 1
