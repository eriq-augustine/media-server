package api;

// This file is to handle community stream specific requests.

import (
   "fmt"
   "net/http"

   "com/pollr/server/database"
   "com/pollr/server/messages"
   "com/pollr/server/static"
);

// Upload a picture.
func uploadImage(userId UserId, eventId int, caption string, imageBase64 string) (interface{}, error) {
   // Note that we are only expecting jpegs.
   path, err := static.UploadBase64Image(imageBase64, "stream", "jpg", fmt.Sprintf("%d", int(userId)));
   if (err != nil) {
      return "", err;
   }

   ok, err := database.AppendToEventSteam(int(userId), eventId, caption, path);
   if (err != nil) {
      return "", err;
   }

   ok, err = database.AppendToUserSteam(int(userId), caption, path);
   if (err != nil) {
      return "", err;
   }

   return messages.NewGeneralStatus(ok, http.StatusOK), nil;
}

func getCommunityStream(eventId int) (interface{}, error) {
   stream, err := database.GetEventStream(eventId);
   if (err != nil) {
      return "", err;
   }

   return messages.NewStream(stream), nil;
}
