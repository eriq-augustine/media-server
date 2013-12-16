#!/bin/sh

echo 'DELETE FROM mediaserver_encodequeue;' | sqlite3 db.sqlite3
echo 'DELETE FROM mediaserver_cache;' | sqlite3 db.sqlite3

rm -f cache/encode/*
rm -f cache/temp/*
