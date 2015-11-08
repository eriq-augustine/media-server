package api;

import (
   "os"

   "com/eriq-augustine/mediaserver/cache"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/messages"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
);

func browsePath(path string) (interface{}, error) {
   // TEST
   log.Debug("Serving browse: " + path);

   path, err := util.RealPath(path);
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
      files = append(files, model.DirEntryFromInfo(fileInfo, path));
   }

   return messages.NewListDir(files), nil;
}

func serveFile(osFile *os.File, path string) (interface{}, error) {
   // TEST
   log.Debug("Serving File: " + path);

   fileInfo, err := osFile.Stat();
   if (err != nil) {
      return "", err;
   }

   file, err := model.NewFile(path, model.DirEntryFromInfo(fileInfo, path));
   if (err != nil) {
      return "", err;
   }

   cache.NegotiateCache(file);

   return messages.NewViewFile(file), nil;
}
