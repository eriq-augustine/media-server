package messages;

import (
   "com/eriq-augustine/mediaserver/model"
)

type ListDir struct {
   Success bool
   IsDir bool
   DirEntries []model.DirEntry
}

func NewListDir(dirEntries []model.DirEntry) *ListDir {
   safeDirEntries := make([]model.DirEntry, 0, len(dirEntries));
   for _, dirEntry := range(dirEntries) {
      safeDirEntries = append(safeDirEntries, dirEntry.Safe());
   }

   return &ListDir{true, true, safeDirEntries};
}

type ViewFile struct {
   Success bool
   IsDir bool
   CacheReady bool
   File model.File
}

func NewViewFile(file model.File, cacheReady bool) *ViewFile {
   return &ViewFile{true, false, cacheReady, file.Safe()};
}
