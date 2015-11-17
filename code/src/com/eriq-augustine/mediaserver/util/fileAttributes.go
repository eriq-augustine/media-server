// +build !windows

package util;

import (
   "os"
   "strings"
)

// Check if a file is hidden.
func IsHidden(fileInfo os.FileInfo) bool {
   // Unix style
   if (strings.HasPrefix(fileInfo.Name(), ".")) {
      return true;
   }

   return false;
}
