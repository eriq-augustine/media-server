#!/bin/sh

if [ -d cache/temp ]; then
   hashes=`ls cache/temp | sed 's/\.[^\.]\+$//'`
   for hash in $hashes ; do
      echo "DELETE FROM mediaserver_cache WHERE hash = '${hash}';" | sqlite3 django_project/db.sqlite3
      rm -f cache/*/${hash}.*
   done
fi
