package server;

import (
   "encoding/hex"
   "flag"
   "fmt"
   "net/http"
   "strings"

   // The client is actually the elfs-api client.
   "github.com/eriq-augustine/elfs-api/apierrors"
   "github.com/eriq-augustine/elfs-api/client"
   "github.com/eriq-augustine/elfs-api/messages"

   "com/eriq-augustine/mediaserver/api"
   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/util"
);

const (
   DEFAULT_BASE_CONFIG_PATH = "config/config.json"
)

// Flags
var (
   configPath = flag.String("config", DEFAULT_BASE_CONFIG_PATH, "Path to the configuration file to use")
   prod = flag.Bool("prod", false, "Use prodution configuration")
)

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

func redirectToHttps(response http.ResponseWriter, request *http.Request) {
   http.Redirect(response, request, fmt.Sprintf("https://%s:%d/%s", request.Host, config.GetInt("httpsPort"), request.RequestURI), http.StatusFound);
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
         jsonResponse, _ := util.ToJSON(messages.NewRejectedToken(apierrors.TokenValidationError{apierrors.TOKEN_VALIDATION_NO_TOKEN}));
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
   clientPrefix := "/" + config.GetString("clientBaseURL") + "/";

   router := api.CreateRouter(clientPrefix);

   // Attach additional prefixes for serving raw and client files.
   // The raw server get auth.
   http.HandleFunc(rawPrefix, AuthFileServer(rawPrefix, config.GetString("staticBaseDir")));

   http.Handle(clientPrefix, BasicFileServer(clientPrefix, config.GetString("clientBaseDir")));
   client.Init();

   http.HandleFunc("/favicon.ico", serveFavicon);
   http.HandleFunc("/robots.txt", serveRobots);

   http.Handle("/", router);

   if (config.GetBool("useSSL")) {
      httpsPort := config.GetInt("httpsPort");

      // Forward http
      if (config.GetBoolDefault("forwardHttp", false) && config.Has("httpPort")) {
         httpPort := config.GetInt("httpPort");

         go func() {
            err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), http.HandlerFunc(redirectToHttps));
            if err != nil {
               log.PanicE("Failed to redirect http to https", err);
            }
         }()
      }

      // Serve https
      log.Info(fmt.Sprintf("Starting media server on https port %d", httpsPort));

      err := http.ListenAndServeTLS(fmt.Sprintf(":%d", httpsPort), config.GetString("httpsCertFile"), config.GetString("httpsKeyFile"), nil);
      if err != nil {
         log.PanicE("Failed to server https", err);
      }
   } else {
      port := config.GetInt("httpPort");
      log.Info(fmt.Sprintf("Starting media server on http port %d", port));

      err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil);
      if err != nil {
         log.PanicE("Failed to server http", err);
      }
   }
}

func LoadConfig() {
   flag.Parse();

   config.LoadFile(*configPath);

   if (*prod) {
      log.SetDebug(false);

      if (config.Has("prodConfig")) {
         config.LoadFile(config.GetString("prodConfig"));
      }
   } else {
      log.SetDebug(true);
   }
}
