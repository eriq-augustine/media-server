package model;

import (
   "os"
   "time"

   "com/eriq-augustine/mediaserver/util"
)

type DirEntry struct {
   Name string
   Size int64
   IsDir bool
   ModTime time.Time
}

func DirEntryFromInfo(fileInfo os.FileInfo) DirEntry {
   return DirEntry{fileInfo.Name(), fileInfo.Size(), fileInfo.IsDir(), fileInfo.ModTime()};
}

type File struct {
   RawLink string
   CacheLink *string
   DirEntry DirEntry
}

func NewFile(path string, dirEnt DirEntry) (File, error) {
   var file File;
   file.DirEntry = dirEnt;

   rawLink, err := util.RawLink(path);
   if (err != nil) {
      return file, err;
   }
   file.RawLink = rawLink;

   ok, cacheLink := util.CacheLink(path)
   if (ok) {
      file.CacheLink = &cacheLink;
   }

   return file, nil;
}
