package server;

import (
   "encoding/hex"
   "fmt"
   "net/http"
   "strings"

   ws "golang.org/x/net/websocket"

   "com/eriq-augustine/mediaserver/api"
   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/messages"
   "com/eriq-augustine/mediaserver/util"
   "com/eriq-augustine/mediaserver/util/errors"
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

func AuthFileServer(urlPrefix string, baseDir string) func(response http.ResponseWriter, request *http.Request) {
   baseHandler := BasicFileServer(urlPrefix, baseDir);

   return func(response http.ResponseWriter, request *http.Request) {
      err := request.ParseForm();
      if (err != nil) {
         log.ErrorE("Unable to parse http request", err);
         jsonResponse, _ := util.ToJSON(messages.NewGeneralStatus(false, http.StatusBadRequest));
         response.WriteHeader(http.StatusBadRequest);
         fmt.Fprintln(response, jsonResponse);
         return;
      }


      token := strings.TrimSpace(request.FormValue("token"));
      _, err = auth.ValidateToken(token);
      if (err != nil) {
         log.WarnE("Bad token request", err);
         jsonResponse, _ := util.ToJSON(messages.NewRejectedToken(errors.TokenValidationError{errors.TOKEN_VALIDATION_NO_TOKEN}));
         response.WriteHeader(http.StatusUnauthorized)
         fmt.Fprintln(response, jsonResponse);
         return;
      }

      baseHandler.ServeHTTP(response, request);
   };
}

func BasicFileServer(urlPrefix string, baseDir string) http.Handler {
   return http.StripPrefix(urlPrefix, http.FileServer(http.Dir(baseDir)));
}

// Note that this will block until the server crashes.
func StartServer() {
   rawPrefix := "/" + config.GetString("rawBaseURL") + "/";
   cachePrefix := "/" + config.GetString("cacheBaseURL") + "/";
   clientPrefix := "/" + config.GetString("clientBaseURL") + "/";

   router := api.CreateRouter(clientPrefix);

   // Attach additional prefixes for serving raw, cached, and client files.
   // The raw and cached servers get auth.
   http.HandleFunc(rawPrefix, AuthFileServer(rawPrefix, config.GetString("staticBaseDir")));
   http.HandleFunc(cachePrefix, AuthFileServer(cachePrefix, config.GetString("cacheBaseDir")));
   http.Handle(clientPrefix, BasicFileServer(clientPrefix, config.GetString("clientBaseDir")));

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

func LoadConfig() {
   config.LoadFile(DEFAULT_BASE_CONFIG_PATH);
   config.LoadFile(DEFAULT_BASE_CONFIG_DEPLOY);
   config.LoadFile(DEFAULT_FILETYPES_CONFIG_PATH);
}
