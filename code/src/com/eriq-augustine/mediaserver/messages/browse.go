package messages;

import (
   "com/eriq-augustine/mediaserver/model"
)

type ListDir struct {
   Success bool
   DirEntries []model.DirEntry
}

func NewListDir(dirEntries []model.DirEntry) *ListDir {
   return &ListDir{true, dirEntries};
}

type ViewFile struct {
   Success bool
   DirEntry model.DirEntry
}

func NewViewFile(dirEntry model.DirEntry) *ViewFile {
   return &ViewFile{true, dirEntry};
}
