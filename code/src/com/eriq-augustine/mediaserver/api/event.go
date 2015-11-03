package api;

import (
   "net/http"

   "com/pollr/server/database"
   "com/pollr/server/event"
   "com/pollr/server/log"
   "com/pollr/server/messages"
);

const DEFAULT_EVENT_COUNT = 5;

func joinEvent(userId UserId, eventId int, accessKey string) interface{} {
   // Check for an accesskey.
   var ok bool = false;
   if (accessKey != "") {
      ok = event.JoinPrivateEvent(int(userId), eventId, accessKey);
   } else {
      ok = event.JoinPublicEvent(int(userId), eventId);
   }

   return messages.NewGeneralStatus(ok, http.StatusOK);
}

func generateEventToken(userId UserId, eventId int, tokenCount int) interface{} {
   // If defaulted (or wrong), then just generate one token.
   if (tokenCount < 1) {
      tokenCount = 1;
   }

   ok, tokens := event.GenerateEventTokens(int(userId), eventId, tokenCount);
   return messages.NewEventTokens(ok, tokens);
}

func getEvents(userId UserId, count int) (interface{}, error) {
   if (count <= 0) {
      count = DEFAULT_EVENT_COUNT;
   }

   events, err := database.GetEvents(int(userId), count);
   if (err != nil) {
      log.ErrorE("Unable to get events", err);
      return "", err;
   }

   return messages.NewEvents(events), nil;
}
