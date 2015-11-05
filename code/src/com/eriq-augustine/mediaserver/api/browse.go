package api;

import (
   "fmt"
   "os"
   "path/filepath"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/messages"
   "com/eriq-augustine/mediaserver/model"
);

func RealStaticPath(path string) (string, error) {
   cleanPath := filepath.Join(config.GetString("staticBaseDir"), strings.TrimPrefix(path, "/"));

   cleanPath, err := filepath.Abs(cleanPath);
   if (err != nil) {
      return "", err;
   }

   cleanPath = filepath.Clean(cleanPath);

   // Ensure that the path is inside the root directory.
   relPath, err := filepath.Rel(config.GetString("staticBaseDir"), cleanPath);
   if (err != nil) {
      return "", err;
   }

   if (strings.HasPrefix(relPath, "..")) {
      return "", fmt.Errorf("Path outside of root");
   }

   return cleanPath, nil;
}

func browsePath(path string) (interface{}, error) {
   // TEST
   log.Debug("Serving browse: " + path);

   path, err := RealStaticPath(path);
   if (err != nil) {
      return "", err;
   }

   // TODO(eriq): First see if file exists and 404 if not.

   file, err := os.Open(path);
   if (err != nil) {
      return "", err;
   }
   defer file.Close();

   fileInfo, err := file.Stat();
   if (err != nil) {
      return "", err;
   } else if (fileInfo.IsDir()) {
      return serveDir(file, path);
   } else {
      return serveFile(file, path);
   }
}

func serveDir(file *os.File, path string) (interface{}, error) {
   // TEST
   log.Debug("Serving Dir: " + path);

   fileInfos, err := file.Readdir(0);
   if (err != nil) {
      return "", err;
   }

   files := make([]model.DirEntry, 0);
   for _, fileInfo := range(fileInfos) {
      files = append(files, model.DirEntryFromInfo(fileInfo));
   }

   return messages.NewListDir(files), nil;
}

func serveFile(file *os.File, path string) (interface{}, error) {
   // TEST
   log.Debug("Serving File: " + path);

   fileInfo, err := file.Stat();
   if (err != nil) {
      return "", err;
   }

   return messages.NewViewFile(model.DirEntryFromInfo(fileInfo)), nil;
}
