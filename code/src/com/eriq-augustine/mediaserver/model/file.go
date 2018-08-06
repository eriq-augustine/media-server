package model;

import (
   "path/filepath"
   "os"
   "time"

   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/util"
)

type DirEntry struct {
   Name string
   Path string
   AbstractPath string
   RawLink *string
   Size int64
   IsDir bool
   ModTime time.Time
}

func DirEntryFromInfo(fileInfo os.FileInfo, path string) DirEntry {
   abstractPath, err := util.AbstractPath(path, "");
   if (err != nil) {
      log.WarnE("Could not get abstract path for (" + path + ")", err);
      abstractPath = path;
   }

   // Files get to show their raw link.
   var rawLinkPtr *string = nil;
   if (!fileInfo.IsDir()) {
      rawLink, err := util.RawLink(filepath.Join(path, fileInfo.Name()));
      if (err != nil) {
         log.WarnE("Could not get raw link for (" + path + ")", err);
      } else {
         rawLinkPtr = &rawLink;
      }
   }

   return DirEntry{fileInfo.Name(), path, abstractPath, rawLinkPtr, fileInfo.Size(), fileInfo.IsDir(), fileInfo.ModTime()};
}

func (dirEntry DirEntry) Safe() DirEntry {
   return DirEntry {
      Name: dirEntry.Name,
      Path: "",
      AbstractPath: dirEntry.AbstractPath,
      RawLink: dirEntry.RawLink,
      Size: dirEntry.Size,
      IsDir: dirEntry.IsDir,
      ModTime: dirEntry.ModTime,
   };
}

type File struct {
   RawLink string
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

func (file File) Safe() File {
   return File{
      RawLink: file.RawLink,
      Poster: file.Poster,
      Subtitles: file.Subtitles,
      DirEntry: file.DirEntry.Safe(),
   };
}
