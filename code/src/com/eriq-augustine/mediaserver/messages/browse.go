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
   File model.File
}

func NewViewFile(file model.File) *ViewFile {
   return &ViewFile{true, file};
}
