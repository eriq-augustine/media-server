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
   PARAM_CAPTION = "caption"
   PARAM_COUNT = "count"
   PARAM_EVENT_ACCESS_TOKEN = "eventAccessToken"
   PARAM_EVENT_ACCESS_TOKEN_COUNT = "tokenCount"
   PARAM_EVENT_ID = "eventId"
   PARAM_FILE = "file"
   PARAM_IMAGE = "image"
   PARAM_PASSHASH = "passhash"
   PARAM_POLL_ID = "pollId"
   PARAM_POLL_ITEM_ID = "pollItemId"
   PARAM_PROFILE = "profile"
   PARAM_PROFILE_USER_ID = "profileUserId"
   PARAM_SEARCH_TEXT = "searchText"
   PARAM_SONG_ID = "songId"
   PARAM_TOKEN = "token"
   PARAM_USERNAME = "username"
   PARAM_USER_ID = "userId"
)

func CreateRouter() *mux.Router {
   methods := []ApiMethod{
      {
         "artist/get",
         getArtists,
         true,
         []ApiMethodParam{
            {PARAM_COUNT, API_PARAM_TYPE_INT, false},
         },
      },
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
      {
         "auth/user/create",
         createAccount,
         false,
         []ApiMethodParam{
            {PARAM_PASSHASH, API_PARAM_TYPE_STRING, true},
            {PARAM_PROFILE, API_PARAM_TYPE_STRING, true},
            {PARAM_IMAGE, API_PARAM_TYPE_STRING, false},
         },
      },
      {
         "community/stream/upload/image",
         uploadImage,
         true,
         []ApiMethodParam{
            {PARAM_EVENT_ID, API_PARAM_TYPE_INT, true},
            {PARAM_CAPTION, API_PARAM_TYPE_STRING, false},
            {PARAM_FILE, API_PARAM_TYPE_STRING, true},
         },
      },
      {
         "community/stream/get",
         getCommunityStream,
         true,
         []ApiMethodParam{
            {PARAM_EVENT_ID, API_PARAM_TYPE_INT, true},
         },
      },
      {
         "event/get/all",
         getEvents,
         true,
         []ApiMethodParam{
            {PARAM_COUNT, API_PARAM_TYPE_INT, false},
         },
      },
      {
         "event/join",
         joinEvent,
         true,
         []ApiMethodParam{
            {PARAM_EVENT_ID, API_PARAM_TYPE_INT, true},
            {PARAM_EVENT_ACCESS_TOKEN, API_PARAM_TYPE_STRING, false},
         },
      },
      {
         "event/generateToken",
         generateEventToken,
         true,
         []ApiMethodParam{
            {PARAM_EVENT_ID, API_PARAM_TYPE_INT, true},
            {PARAM_EVENT_ACCESS_TOKEN_COUNT, API_PARAM_TYPE_INT, false},
         },
      },
      {
         "home/get/recent",
         getRecentHome,
         true,
         []ApiMethodParam{
            {PARAM_COUNT, API_PARAM_TYPE_INT, false},
         },
      },
      {
         "news/get/recent",
         getRecentNews,
         true,
         []ApiMethodParam{
            {PARAM_COUNT, API_PARAM_TYPE_INT, false},
         },
      },
      {
         "poll/get",
         getPoll,
         true,
         []ApiMethodParam{
            {PARAM_POLL_ID, API_PARAM_TYPE_INT, true},
         },
      },
      {
         "poll/vote",
         vote,
         true,
         []ApiMethodParam{
            {PARAM_POLL_ITEM_ID, API_PARAM_TYPE_INT, true},
         },
      },
      {
         "poll/request",
         request,
         true,
         []ApiMethodParam{
            {PARAM_POLL_ID, API_PARAM_TYPE_INT, true},
            {PARAM_SONG_ID, API_PARAM_TYPE_INT, true},
         },
      },
      {
         "profile/get",
         getProfile,
         true,
         []ApiMethodParam{
            {PARAM_PROFILE_USER_ID, API_PARAM_TYPE_INT, true},
         },
      },
      {
         "profile/set",
         setProfile,
         true,
         []ApiMethodParam{
            {PARAM_PROFILE, API_PARAM_TYPE_STRING, true},
            {PARAM_IMAGE, API_PARAM_TYPE_STRING, false},
         },
      },
      {
         "search/songs",
         searchSongs,
         true,
         []ApiMethodParam{
            {PARAM_SEARCH_TEXT, API_PARAM_TYPE_STRING, true},
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

   return router;
}
