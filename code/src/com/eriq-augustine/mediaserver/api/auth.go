package api;

// Implements the "/auth/" portion of the api.
// (This file is not about implementing the authentication middleware.)

import (
   "net/http"

   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/messages"
   "com/eriq-augustine/mediaserver/util/errors"
);

// Invalidating a token is akin to logging out.
// Note that one must have a valid token to invalidate their own token.
func invalidateToken(token Token) (interface{}, error) {
   ok, err := auth.InvalidateToken(string(token));

   if (err != nil) {
      return "", err;
   }

   return messages.NewGeneralStatus(ok, http.StatusOK), nil;
}

func requestToken(username string, passhash string) (interface{}, int, error) {
   // userId, token, err := auth.AuthenticateUser(username, passhash);
   _, token, err := auth.AuthenticateUser(username, passhash);
   if (err != nil) {
      validationErr, ok := err.(errors.TokenValidationError);
      if (!ok) {
         // Some other (non-validation) error.
         return "", 0, err;
      } else {
         return messages.NewRejectedToken(validationErr), http.StatusForbidden, err;
      }
   } else {
      return messages.NewAuthorizedToken(token), 0, nil;
   }
}

/*
func createAccount(passhash string, stringProfile string) (interface{}, int, error) {
}
*/
