#!/bin/sh

echo 'DELETE FROM mediaserver_encodequeue;' | sqlite3 django_project/db.sqlite3
echo 'DELETE FROM mediaserver_cache;' | sqlite3 django_project/db.sqlite3

rm -Rf cache/encode/*
rm -Rf cache/temp/*
