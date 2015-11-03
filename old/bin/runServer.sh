#!/bin/sh

bin/start_encode.sh
bin/start_websocket.sh
bin/gunicorn_start.sh
