package api;

// Implements the "/search/" portion of the api.

import (
   "com/pollr/server/database"
   "com/pollr/server/log"
   "com/pollr/server/messages"
);

func searchSongs(searchText string) (interface{}, error) {
   songs, err := database.SearchSongs(searchText);
   if (err != nil) {
      log.ErrorE("Unable to search songs", err);
      return "", err;
   }

   return messages.NewSearchResults(songs), nil;
}
