package util;

import (
   "path/filepath"
   "os"

   "github.com/eriq-augustine/golog"
)

func RmDir(path string) {
   os.RemoveAll(path);
}

func DirSize(path string) uint64 {
   var size uint64 = 0;

   err := filepath.Walk(path, func(childPath string, fileInfo os.FileInfo, err error) error {
      if (err != nil) {
         return err;
      }

      // Skip root.
      if (path == childPath) {
         return nil;
      }

      if (fileInfo.IsDir()) {
         size += DirSize(childPath);
      } else {
         size += uint64(fileInfo.Size());
      }

      return nil;
   });

   if (err != nil) {
      golog.ErrorE("Error getting the a directory's size (" + path + ")", err);
      return 0;
   }

   return size;
}
