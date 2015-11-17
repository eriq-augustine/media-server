package util;

import (
   "os"
   "strings"
   "syscall"
)

// Check if a file is hidden.
func IsHidden(fileInfo os.FileInfo) bool {
   // Check for windows hidden attribute.
   sysInfo := fileInfo.Sys();
   if (sysInfo != nil) {
      winAttributes, ok := sysInfo.(syscall.Win32FileAttributeData);
      if (ok) {
         if (winAttributes.FileAttributes & syscall.FILE_ATTRIBUTE_HIDDEN != 0) {
            return true;
         }
      }
   }

   return false;
}
