package main;

import (
   "encoding/hex"
   "fmt"
   "net/http"

   ws "golang.org/x/net/websocket"

   "com/eriq-augustine/mediaserver/api"
   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/websocket"
);

const (
   DEFAULT_BASE_CONFIG_PATH = "config/config-base.json"
   DEFAULT_BASE_CONFIG_DEPLOY = "config/config-deploy.json"
   DEFAULT_FILETYPES_CONFIG_PATH = "config/filetypes.json"
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

func redirectToClient(response http.ResponseWriter, request *http.Request) {
   http.Redirect(response, request, "http://localhost:1234/client/", 200);
}

func main() {
   config.LoadFile(DEFAULT_BASE_CONFIG_PATH);
   config.LoadFile(DEFAULT_BASE_CONFIG_DEPLOY);
   config.LoadFile(DEFAULT_FILETYPES_CONFIG_PATH);

   // It is safe to load users after the configs have been loaded.
   auth.LoadUsers();

   router := api.CreateRouter("/" + config.GetString("clientBaseURL") + "/");

   // Attach additional prefixes for serving raw and cached files.
   rawPrefix := "/" + config.GetString("rawBaseURL") + "/";
   http.Handle(rawPrefix, http.StripPrefix(rawPrefix, http.FileServer(http.Dir(config.GetString("staticBaseDir")))));

   cachePrefix := "/" + config.GetString("cacheBaseURL") + "/";
   http.Handle(cachePrefix, http.StripPrefix(cachePrefix, http.FileServer(http.Dir(config.GetString("cacheBaseDir")))));

   clientPrefix := "/" + config.GetString("clientBaseURL") + "/";
   http.Handle(clientPrefix, http.StripPrefix(clientPrefix, http.FileServer(http.Dir(config.GetString("clientBaseDir")))));

   http.HandleFunc("/favicon.ico", serveFavicon);
   http.HandleFunc("/robots.txt", serveRobots);

   // Websocket
   http.Handle("/ws", ws.Handler(websocket.SocketHandler));

   http.Handle("/", router);

   port := config.GetIntDefault("port", 1234);
   log.Info(fmt.Sprintf("Starting media server on port %d", port));

   err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil);
   if err != nil {
      panic("ListenAndServe: " + err.Error());
   }
}
