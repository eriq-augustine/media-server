package api;

// The definition of the methods used for the API.
// See apimethod.go before trying to make any new methods.

import (
   "net/http"

   "github.com/gorilla/mux"
)

// We need to define these as types, so we can figure it out when we want to pass params.
type Token string;
type UserId int;

const (
   PARAM_FILE = "file"
   PARAM_IMAGE = "image"
   PARAM_PASSHASH = "passhash"
   PARAM_PATH = "path"
   PARAM_TOKEN = "token"
   PARAM_USERNAME = "username"
   PARAM_USER_ID = "userId"
)

func CreateRouter(rootRedirect string) *mux.Router {
   methods := []ApiMethod{
      {
         "auth/token/request",
         requestToken,
         false,
         []ApiMethodParam{
            {PARAM_USERNAME, API_PARAM_TYPE_STRING, true},
            {PARAM_PASSHASH, API_PARAM_TYPE_STRING, true},
         },
      },
      {
         "auth/token/invalidate",
         invalidateToken,
         true,
         []ApiMethodParam{},
      },
      /*
      {
         "auth/user/create",
         createAccount,
         false,
         []ApiMethodParam{
            {PARAM_PASSHASH, API_PARAM_TYPE_STRING, true},
            {PARAM_PROFILE, API_PARAM_TYPE_STRING, true},
         },
      },
      */
      {
         "browse/path",
         browsePath,
         false, // TODO(eriq): Auth
         []ApiMethodParam{
            {PARAM_PATH, API_PARAM_TYPE_STRING, false},
         },
      },
      {
         "serve/path",
         servePath,
         false, // TODO(eriq): Auth
         []ApiMethodParam{
            {PARAM_PATH, API_PARAM_TYPE_STRING, false},
         },
      },
   };

   router := mux.NewRouter();
   for _, method := range(methods) {
      method.Validate();
      router.HandleFunc(buildApiUrl(method.Path), method.Middleware());
   }

   // Handle 404 specially.
   var notFoundApiMethod ApiMethod = ApiMethod{
      "__404__", // We will not actually bind 404 to a path, so just use something to pass validation.
      notFound,
      true, // We don't give hints about our API, so require auth for everything.
      []ApiMethodParam{}, // Not expecting any params for 404.
   };
   notFoundApiMethod.Validate();
   router.NotFoundHandler = http.HandlerFunc(notFoundApiMethod.Middleware());

   // If supplied, register the root redirect.
   // Root should never be hit directly, so we can optionally redirect it.
   if (rootRedirect != "") {
      router.Handle("/", http.RedirectHandler(rootRedirect, 301));
   }

   return router;
}
