package model;

import (
   "os"
)

type DirEntry struct {
   Name string
   Size int64
   IsDir bool
}

func DirEntryFromInfo(fileInfo os.FileInfo) DirEntry {
   return DirEntry{fileInfo.Name(), fileInfo.Size(), fileInfo.IsDir()};
}
