package api;

// Implements the "/artist/" portion of the api.

import (
   "com/pollr/server/database"
   "com/pollr/server/log"
   "com/pollr/server/messages"
);

func getArtists(count int) (interface{}, error) {
   if (count == 0) {
      count = 26;
   }

   artists, err := database.GetArtists(count);
   if (err != nil) {
      log.ErrorE("Unable to get artists", err);
      return "", err;
   }

   return messages.NewArtistResult(artists), nil;
}
