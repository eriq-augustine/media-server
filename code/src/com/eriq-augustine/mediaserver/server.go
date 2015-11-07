package main;

import (
   "encoding/hex"
   "flag"
   "fmt"
   "net/http"

   "com/eriq-augustine/mediaserver/api"
   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
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

   router := api.CreateRouter();

   // Attach additional prefixes for serving raw and cached files.
   rawPrefix := "/" + config.GetString("rawBaseURL") + "/";
   http.Handle(rawPrefix, http.StripPrefix(rawPrefix, http.FileServer(http.Dir(config.GetString("staticBaseDir")))));

   cachePrefix := "/" + config.GetString("cacheBaseURL") + "/";
   http.Handle(cachePrefix, http.StripPrefix(cachePrefix, http.FileServer(http.Dir(config.GetString("cacheBaseDir")))));

   clientPrefix := "/" + config.GetString("clientBaseURL") + "/";
   http.Handle(clientPrefix, http.StripPrefix(clientPrefix, http.FileServer(http.Dir(config.GetString("clientBaseDir")))));

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
