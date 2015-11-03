package api;

// Implements the "/auth/" portion of the api.
// (This file is not about implementing the authentication middleware.)

import (
   "net/http"
   "sync"

   "com/pollr/server/database"
   "com/pollr/server/auth"
   "com/pollr/server/log"
   "com/pollr/server/messages"
   "com/pollr/server/model"
   "com/pollr/server/static"
   "com/pollr/server/util/errors"
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
   userId, token, err := auth.AuthenticateUser(username, passhash);
   if (err != nil) {
      validationErr, ok := err.(errors.TokenValidationError);
      if (!ok) {
         // Some other (non-validation) error.
         return "", 0, err;
      } else {
         return messages.NewRejectedToken(validationErr), http.StatusForbidden, err;
      }
   } else {
      user, err := database.GetOwnUserProfile(userId);
      if (err != nil) {
         return "", 0, err;
      }

      return messages.NewAuthorizedToken(token, *user), 0, nil;
   }
}

func createAccount(passhash string, stringProfile string, imageBase64 string) (interface{}, int, error) {
   // Try to convert the profile.
   profile, err := model.JSONToProfile(stringProfile);

   if (err != nil) {
      log.WarnE("Bad profile JSON for account creation", err);
      // Error doesn't matter here.
      return messages.NewGeneralStatus(false, http.StatusBadRequest), http.StatusBadRequest, nil;
   }

   err = model.ValidateProfile(profile);
   if (err != nil) {
      log.WarnE("Invalid profile for account creation", err);
      return messages.NewGeneralStatus(false, http.StatusBadRequest), http.StatusBadRequest, nil;
   }

   if (imageBase64 != "") {
      // Just use a pretty much random seed (c for "Create Account") since we don't have a user id.
      // We don't technically need a seed, but it won't hurt.
      path, err := static.UploadBase64Image(imageBase64, "user", "jpg", "c");
      if (err != nil) {
         log.WarnE("Invalid profile image", err);
         return messages.NewGeneralStatus(false, http.StatusBadRequest), http.StatusBadRequest, nil;
      }

      profile.ProfileImagePath = &path;
   }

   // Put a lock around user creation.
   mutex := &sync.Mutex{}
   mutex.Lock();

   // Unlock once we are all done.
   defer mutex.Unlock();

   // Ensure that the username and email are not already taken.
   userExists, err := database.UserExists(*profile.Username);
   if (err != nil) {
      log.ErrorE("Unable to check if user exists for account creation", err);
      return "", 0, err;
   }

   emailExists, err := database.UserEmailExists(*profile.Email);
   if (err != nil) {
      log.ErrorE("Unable to check if user email exists for account creation", err);
      return "", 0, err;
   }

   if (userExists || emailExists) {
      return messages.NewAccountCreation(userExists, emailExists, "", nil), 0, nil;
   }

   // Create the actual account.
   _, err = database.CreateUser(profile, passhash);
   if (err != nil) {
      log.ErrorE("Unable to create user", err);
      return "", 0, err;
   }

   userId, token, _ := auth.AuthenticateUser(*profile.Username, passhash);
   user, err := database.GetOwnUserProfile(userId);
   if (err != nil) {
      return "", 0, err;
   }

   return messages.NewAccountCreation(false, false, token, user), 0, nil;
}
