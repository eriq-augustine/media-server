package main;

import (
   "encoding/hex"
   "flag"
   "fmt"
   "net/http"
   "os"
   "path/filepath"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/util"

   "github.com/gorilla/mux"
);

const (
   DEFAULT_CONFIG_PATH = "config/config.json"
   DEFAULT_DEV_CONFIG_PATH = "config/config-dev.json"
   DEFAULT_PROD_CONFIG_PATH = "config/config-prod.json"
   DEFAULT_SECRETS_PATH = "config/secrets.json"
);

func serveFavicon(response http.ResponseWriter, request *http.Request) {
   dataBytes, err := hex.DecodeString(config.GetStringDefault("favicon", ""));

   if (err != nil) {
      response.WriteHeader(http.StatusInternalServerError);
      return;
   }

   response.WriteHeader(http.StatusOK);
   response.Header().Set("Content-Type", "image/x-icon");
   response.Write(dataBytes);
}

func serveRobots(response http.ResponseWriter, request *http.Request) {
   fmt.Fprintf(response, "User-agent: *\nDisallow: /\n");
}

func serveStatic(response http.ResponseWriter, request *http.Request, fileServer http.Handler) {
   vars := mux.Vars(request)
   path := vars["path"]

   log.Debug("Serving static: " + path);

   fileServer.ServeHTTP(response, request);
}

func RealStaticPath(path string) string {
   cleanPath := filepath.Join(config.GetString("staticBaseDir"), strings.TrimPrefix(path, "/"));

   // TODO(eriq): Better handling
   cleanPath, err := filepath.Abs(cleanPath);
   if (err != nil) {
      log.PanicE("Bad static path: " + cleanPath, err);
   }

   cleanPath = filepath.Clean(cleanPath);

   // Ensure that the path is inside the root directory.
   relPath, err := filepath.Rel(config.GetString("staticBaseDir"), cleanPath);
   // TODO(eriq): Better error handling
   if (err != nil) {
      log.PanicE("Cannot get relative path", err);
   }

   if (strings.HasPrefix(relPath, "..")) {
      log.Panic("Path not inside root");
   }

   return cleanPath;
}

type File struct {
   Name string
   Size int64
   IsDir bool
}

func FileFromInfo(fileInfo os.FileInfo) File {
   return File{fileInfo.Name(), fileInfo.Size(), fileInfo.IsDir()};
}

func serveBrowse(response http.ResponseWriter, request *http.Request) {
   vars := mux.Vars(request)
   path := vars["path"]

   // TEST
   fmt.Println("Serving browse: " + path);

   path = RealStaticPath(path);

   file, err := os.Open(path);
   if (err != nil) {
      // TODO(eriq): Return 404
      log.ErrorE("Could not get file", err);
      fmt.Fprintln(response, "404");
      return;
   }
   defer file.Close();

   fileInfo, err := file.Stat();
   if (err != nil) {
      // TODO(eriq): Return 404
      log.ErrorE("Could not stat file", err);
      fmt.Fprintln(response, "404");
   } else if (fileInfo.IsDir()) {
      serveDir(response, file, path);
   } else {
      serveFile(response, file, path);
   }
}

func serveDir(response http.ResponseWriter, file *os.File, path string) {
   // TEST
   fmt.Println("Serving Dir: " + path);

   fileInfos, err := file.Readdir(0);
   if (err != nil) {
      log.ErrorE("Unable to readdir", err);
      // TODO(eriq): Return better
      return;
   }

   files := make([]File, 0);
   for _, fileInfo := range(fileInfos) {
      files = append(files, FileFromInfo(fileInfo));
   }

   response.Header().Set("Content-Type", "application/json; charset=UTF-8");
   jsonResponse, _ := util.ToJSON(files);
   fmt.Fprintln(response, jsonResponse);
}

func serveFile(response http.ResponseWriter, file *os.File, path string) {
   // TEST
   fmt.Println("Serving File: " + path);

   fileInfo, err := file.Stat();
   if (err != nil) {
      // TODO(eriq): Log better
      log.ErrorE("Can't stat file", err);
      return;
   }

   response.Header().Set("Content-Type", "application/json; charset=UTF-8");
   jsonResponse, _ := util.ToJSON(FileFromInfo(fileInfo));
   fmt.Fprintln(response, jsonResponse);
}

func main() {
   var prod *bool = flag.Bool("prod", false, "Run in production mode");
   flag.Parse();

   config.LoadFile(DEFAULT_CONFIG_PATH);
   config.LoadFile(DEFAULT_SECRETS_PATH);

   if (*prod) {
      config.LoadFile(DEFAULT_PROD_CONFIG_PATH);
   } else {
      log.SetDebug(true);
      config.LoadFile(DEFAULT_DEV_CONFIG_PATH);
   }

   fileServer := http.StripPrefix("/" + config.GetString("staticBaseURL"), http.FileServer(http.Dir(config.GetString("staticBaseDir"))));
   staticHandler := func(response http.ResponseWriter, request *http.Request) {
      serveStatic(response, request, fileServer);
   };

   router := mux.NewRouter();
   router.HandleFunc("/" + config.GetString("staticBaseURL") + "/{path:.*}", staticHandler);
   router.HandleFunc("/" + config.GetString("browseBaseURL") + "/{path:.*}", serveBrowse)

   http.HandleFunc("/favicon.ico", serveFavicon);
   http.HandleFunc("/robots.txt", serveRobots);

   http.Handle("/", router);

   port := config.GetIntDefault("port", 1234);
   log.Info(fmt.Sprintf("Starting media server on port %d", port));

   err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil);
   if err != nil {
      panic("ListenAndServe: " + err.Error());
   }
}
