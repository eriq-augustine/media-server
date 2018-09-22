package api;

import (
   "mime"
   "net/http"
   "os"

   "github.com/eriq-augustine/elfs-api/messages"
   "github.com/eriq-augustine/elfs-api/model"
   "github.com/eriq-augustine/goconfig"
   "github.com/eriq-augustine/golog"
   "github.com/pkg/errors"

   "com/eriq-augustine/mediaserver/util"
);

func init() {
   // Add webvtt into the mime type.
   mime.AddExtensionType(".vtt", "text/vtt");
}

func browsePath(path string) (interface{}, int, error) {
   golog.Debug("Serving: " + path);

   path, err := util.RealPath(path);
   if (err != nil) {
      return "", 0, err;
   }

   file, fileInfo, message, status, err := getFileInfo(path);
   if (file != nil) {
      defer file.Close();
   }

   if (message != nil || status != 0 || err != nil) {
      return message, status, err;
   }

   if (fileInfo.IsDir()) {
      return serveDir(file, fileInfo, path);
   } else {
      return messages.NewFileInfo(util.DirEntryFromInfo(fileInfo, path)), 0, nil;
   }
}

func serveDir(file *os.File, fileInfo os.FileInfo, path string) (interface{}, int, error) {
   children, err := file.Readdir(0);
   if (err != nil) {
      return "", 0, err;
   }

   showHidden := goconfig.GetBoolDefault("showHiddenFiles", false);

   dirents := make([]*model.DirEntry, 0);
   for _, childFileInfo := range(children) {
      if (!showHidden && util.IsHidden(childFileInfo)) {
         continue;
      }

      var childPath string = util.Join(path, childFileInfo.Name());
      dirents = append(dirents, util.DirEntryFromInfo(childFileInfo, childPath));
   }

   return messages.NewListDir(util.DirEntryFromInfo(fileInfo, path), dirents), 0, nil;
}

func getFileContents(path string) (interface{}, int, string, error) {
   golog.Debug("Serving Contents: [" + path + "]");

   path, err := util.RealPath(path);
   if (err != nil) {
      return "", 0, "", err;
   }

   file, fileInfo, message, status, err := getFileInfo(path);
   if (message != nil || status != 0 || err != nil) {
      return message, status, "", err;
   }

   if (fileInfo.IsDir()) {
      file.Close();
      return "", http.StatusBadRequest, "", errors.New("Cannot get the file contents of a dir.");
   }

   return file, 0, mime.TypeByExtension("." + util.Ext(fileInfo.Name())), nil;
}

// Caller must close the file.
func getFileInfo(path string) (*os.File, os.FileInfo, interface{}, int, error) {
   if (!util.PathExists(path)) {
      return nil, nil, messages.NewGeneralStatus(false, http.StatusNotFound), http.StatusNotFound, nil;
   }

   file, err := os.Open(path);
   if (err != nil) {
      return nil, nil, "", 0, err;
   }

   fileInfo, err := file.Stat();
   if (err != nil) {
      return nil, nil, "", 0, err;
   }

   // 404 if we shouldn't be seeing this file.
   if (!goconfig.GetBoolDefault("showHiddenFiles", false) && util.IsHidden(fileInfo)) {
      return nil, nil, messages.NewGeneralStatus(false, http.StatusNotFound), http.StatusNotFound, nil;
   }

   return file, fileInfo, nil, 0, nil;
}
