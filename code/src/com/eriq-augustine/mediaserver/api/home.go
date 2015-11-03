package api;

// Implements the "/home/" portion of the api.

import (
   "com/pollr/server/database"
   "com/pollr/server/log"
   "com/pollr/server/messages"
);

const DEFAULT_HOME_COUNT = 5;

func getRecentHome(count int) (interface{}, error) {
   if (count <= 0) {
      count = DEFAULT_HOME_COUNT;
   }

   homeArticles, err := database.GetRecentHome(count);
   if (err != nil) {
      log.ErrorE("Unable to get recent home articles", err);
      return "", err;
   }

   return messages.NewHomeArticles(homeArticles), nil;
}
