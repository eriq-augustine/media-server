package api;

import (
   "mime"
   "net/http"
   "os"

   "github.com/eriq-augustine/elfs-api/messages"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
);

func init() {
   // Add webvtt into the mime type.
   mime.AddExtensionType(".vtt", "text/vtt");
}

func browsePath(path string) (interface{}, int, error) {
   log.Debug("Serving: " + path);

   path, err := util.RealPath(path);
   if (err != nil) {
      return "", 0, err;
   }

   if (!util.PathExists(path)) {
      return messages.NewGeneralStatus(false, http.StatusNotFound), http.StatusNotFound, nil;
   }

   file, err := os.Open(path);
   if (err != nil) {
      return "", 0, err;
   }
   defer file.Close();

   fileInfo, err := file.Stat();
   if (err != nil) {
      return "", 0, err;
   }

   // 404 if we shouldn't be seeing this file.
   if (!config.GetBoolDefault("showHiddenFiles", false) && util.IsHidden(fileInfo)) {
      return messages.NewGeneralStatus(false, http.StatusNotFound), http.StatusNotFound, nil;
   }

   if (fileInfo.IsDir()) {
      return serveDir(file, path);
   } else {
      return serveFile(file, path);
   }
}

func serveDir(file *os.File, path string) (interface{}, int, error) {
   children, err := file.Readdir(0);
   if (err != nil) {
      return "", 0, err;
   }

   showHidden := config.GetBoolDefault("showHiddenFiles", false);

   dirents := make([]model.DirEntry, 0);
   for _, fileInfo := range(children) {
      if (!showHidden && util.IsHidden(fileInfo)) {
         continue;
      }

      dirents = append(dirents, model.DirEntryFromInfo(fileInfo, path));
   }

   return messages.NewListDir(model.DirEntryFromInfo(fileInfo, path), dirents), 0, nil;
}

func serveFile(osFile *os.File, path string) (interface{}, int, error) {
   fileInfo, err := osFile.Stat();
   if (err != nil) {
      return "", 0, err;
   }

   file, err := model.NewFile(path, model.DirEntryFromInfo(fileInfo, path));
   if (err != nil) {
      return "", 0, err;
   }

   return messages.NewViewFile(file), 0, nil;
}
