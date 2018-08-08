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
   "github.com/eriq-augustine/goconfig"
   "github.com/eriq-augustine/golog"

   "com/eriq-augustine/mediaserver/api"
   "com/eriq-augustine/mediaserver/auth"
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
   dataBytes, err := hex.DecodeString(goconfig.GetStringDefault("favicon", ""));

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
   http.Redirect(response, request, fmt.Sprintf("https://%s:%d/%s", request.Host, goconfig.GetInt("httpsPort"), request.RequestURI), http.StatusFound);
}

func AuthFileServer(urlPrefix string, baseDir string) func(response http.ResponseWriter, request *http.Request) {
   baseHandler := BasicFileServer(urlPrefix, baseDir);

   return func(response http.ResponseWriter, request *http.Request) {
      err := request.ParseForm();
      if (err != nil) {
         golog.ErrorE("Unable to parse http request", err);
         jsonResponse, _ := util.ToJSON(messages.NewGeneralStatus(false, http.StatusBadRequest));
         response.WriteHeader(http.StatusBadRequest);
         fmt.Fprintln(response, jsonResponse);
         return;
      }

      token := strings.TrimSpace(request.FormValue("token"));
      _, err = auth.ValidateToken(token);
      if (err != nil) {
         golog.WarnE("Bad token request", err);
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
   clientPrefix := "/" + goconfig.GetString("clientBaseURL") + "/";

   router := api.CreateRouter(clientPrefix);

   // Serve the client.
   http.Handle(clientPrefix, BasicFileServer(clientPrefix, goconfig.GetString("clientBaseDir")));
   client.Init();

   http.HandleFunc("/favicon.ico", serveFavicon);
   http.HandleFunc("/robots.txt", serveRobots);

   http.Handle("/", router);

   if (goconfig.GetBool("useSSL")) {
      httpsPort := goconfig.GetInt("httpsPort");

      // Forward http
      if (goconfig.GetBoolDefault("forwardHttp", false) && goconfig.Has("httpPort")) {
         httpPort := goconfig.GetInt("httpPort");

         go func() {
            err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), http.HandlerFunc(redirectToHttps));
            if err != nil {
               golog.PanicE("Failed to redirect http to https", err);
            }
         }()
      }

      // Serve https
      golog.Info(fmt.Sprintf("Starting media server on https port %d", httpsPort));

      err := http.ListenAndServeTLS(fmt.Sprintf(":%d", httpsPort), goconfig.GetString("httpsCertFile"), goconfig.GetString("httpsKeyFile"), nil);
      if err != nil {
         golog.PanicE("Failed to server https", err);
      }
   } else {
      port := goconfig.GetInt("httpPort");
      golog.Info(fmt.Sprintf("Starting media server on http port %d", port));

      err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil);
      if err != nil {
         golog.PanicE("Failed to server http", err);
      }
   }
}

func LoadConfig() {
   flag.Parse();

   goconfig.LoadFile(*configPath);

   if (*prod) {
      golog.SetDebug(false);

      if (goconfig.Has("prodConfig")) {
         goconfig.LoadFile(goconfig.GetString("prodConfig"));
      }
   } else {
      golog.SetDebug(true);
   }
}
