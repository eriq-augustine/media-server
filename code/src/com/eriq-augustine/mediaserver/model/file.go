package model;

import (
   "os"
   "time"

   "com/eriq-augustine/mediaserver/util"
)

type DirEntry struct {
   Name string
   Path string
   Size int64
   IsDir bool
   ModTime time.Time
}

func DirEntryFromInfo(fileInfo os.FileInfo, path string) DirEntry {
   return DirEntry{fileInfo.Name(), path, fileInfo.Size(), fileInfo.IsDir(), fileInfo.ModTime()};
}

type File struct {
   RawLink string
   CacheLink *string
   Poster *string
   Subtitles []string
   DirEntry DirEntry
}

func NewFile(path string, dirEnt DirEntry) (File, error) {
   var file File;
   file.DirEntry = dirEnt;
   file.Subtitles = make([]string, 0);

   rawLink, err := util.RawLink(path);
   if (err != nil) {
      return file, err;
   }
   file.RawLink = rawLink;

   return file, nil;
}
