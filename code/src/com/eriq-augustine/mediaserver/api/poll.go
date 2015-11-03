package api;

// Implements the "/poll/" portion of the api.

import (
   "net/http"

   "com/pollr/server/database"
   "com/pollr/server/log"
   "com/pollr/server/messages"
);

// Vote for an existing poll item.
func vote(userId UserId, pollItemId int) (interface{}, error) {
   ok, err := database.VoteInPoll(int(userId), pollItemId);

   if (err != nil) {
      return "", err;
   }

   return messages.NewGeneralStatus(ok, http.StatusOK), nil;
}

// Request a new poll item.
func request(userId UserId, pollId int, songId int) (interface{}, error) {
   ok, err := database.RequestPollItem(int(userId), pollId, songId);

   if (err != nil) {
      return "", err;
   }

   return messages.NewGeneralStatus(ok, http.StatusOK), nil;
}

func getPoll(pollId int, userId UserId) (interface{}, error) {
   pollItems, err := database.GetPoll(pollId, int(userId));
   if (err != nil) {
      log.ErrorE("Unable to get poll", err);
      return "", err;
   }

   return messages.NewPoll(pollItems), nil;
}
