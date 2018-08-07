package api;

// The definition of the methods used for the API.
// See apimethod.go before trying to make any new methods.

import (
   "fmt"
   "net/http"
   "strings"

   "github.com/gorilla/mux"
   "github.com/eriq-augustine/elfs-api/messages"
   "github.com/eriq-augustine/goapi"
   "github.com/eriq-augustine/golog"

   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/config"
)

const (
   PARAM_FILE = "file"
   PARAM_IMAGE = "image"
   PARAM_PASSHASH = "passhash"
   PARAM_ID = "id"
   PARAM_TOKEN = "token"
   PARAM_USERNAME = "username"
)

func validateToken(token string, log goapi.Logger) (int, string, error) {
   userName, err := auth.ValidateToken(token);
   return 0, userName, err;
}

func CreateRouter(rootRedirect string) *mux.Router {
   var factory goapi.ApiMethodFactory;

   factory.SetLogger(golog.Logger{});
   factory.SetTokenValidator(validateToken);

   methods := []*goapi.ApiMethod{
      factory.NewApiMethod(
         "auth/token/request",
         requestToken,
         false,
         []goapi.ApiMethodParam{
            {PARAM_USERNAME, goapi.API_PARAM_TYPE_STRING, true},
            {PARAM_PASSHASH, goapi.API_PARAM_TYPE_STRING, true},
         },
      ),
      factory.NewApiMethod(
         "auth/token/invalidate",
         invalidateToken,
         true,
         []goapi.ApiMethodParam{},
      ),
      factory.NewApiMethod(
         "browse",
         browsePath,
         true,
         []goapi.ApiMethodParam{
            {PARAM_ID, goapi.API_PARAM_TYPE_STRING, false},
         },
      ),
      factory.NewApiMethod(
         "contents",
         getFileContents,
         true,
         []goapi.ApiMethodParam{
            {PARAM_ID, goapi.API_PARAM_TYPE_STRING, true},
         },
      ).SetAllowTokenParam(true),
   };

   router := mux.NewRouter();
   for _, method := range(methods) {
      router.HandleFunc(buildApiUrl(method.Path()), method.Middleware());
   }

   // Handle 404 specially.
   var notFoundApiMethod *goapi.ApiMethod = factory.NewApiMethod(
      "__404__", // We will not actually bind 404 to a path, so just use something to pass validation.
      notFound,
      true, // We don't give hints about our API, so require auth for everything.
      []goapi.ApiMethodParam{}, // Not expecting any params for 404.
   );
   router.NotFoundHandler = http.HandlerFunc(notFoundApiMethod.Middleware());

   // If supplied, register the root redirect.
   // Root should never be hit directly, so we can optionally redirect it.
   if (rootRedirect != "") {
      router.Handle("/", http.RedirectHandler(rootRedirect, 301));
   }

   return router;
}

func buildApiUrl(path string) string {
   path = strings.TrimPrefix(path, "/");

   return fmt.Sprintf("/api/v%02d/%s", config.GetIntDefault("apiVersion", 0), path);
}

func notFound() (interface{}, int) {
   return messages.NewGeneralStatus(false, http.StatusNotFound), http.StatusNotFound;
}
