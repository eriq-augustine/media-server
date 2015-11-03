package api;

// Implements profile related api functionality.

import (
   "fmt"
   "net/http"

   "com/pollr/server/database"
   "com/pollr/server/log"
   "com/pollr/server/messages"
   "com/pollr/server/model"
   "com/pollr/server/static"
);

// Get a user's profile.
// You can get any profile, but the content will vary based on whether it is your own profile.
func getProfile(userId UserId, profileUserId int) (interface{}, error) {
   var profile *model.Profile = nil;
   var err error = nil;
   if (int(userId) == profileUserId) {
      profile, err = database.GetOwnUserProfile(int(userId));
   } else {
      profile, err = database.GetUserProfile(profileUserId);
   }

   if (err != nil) {
      log.ErrorE("Unable to get profile", err);
      return "", err;
   }

   return messages.NewProfile(profile), nil;
}

func setProfile(userId UserId, stringProfile string, imageBase64 string) (interface{}, int, error) {
   // Try to convert the profile.
   profile, err := model.JSONToProfile(stringProfile);

   if (err != nil) {
      log.WarnE("Bad profile JSON", err);
      // Error doesn't matter here.
      return messages.NewGeneralStatus(false, http.StatusBadRequest), http.StatusBadRequest, nil;
   }

   err = model.ValidateProfile(profile);
   if (err != nil) {
      log.WarnE("Invalid profile for set", err);
      return messages.NewGeneralStatus(false, http.StatusBadRequest), http.StatusBadRequest, nil;
   }

   if (imageBase64 != "") {
      path, err := static.UploadBase64Image(imageBase64, "user", "jpg", fmt.Sprintf("%d", int(userId)));
      if (err != nil) {
         log.WarnE("Invalid profile image", err);
         return messages.NewGeneralStatus(false, http.StatusBadRequest), http.StatusBadRequest, nil;
      }

      profile.ProfileImagePath = &path;
   }

   ok, err := database.SetUserProfile(int(userId), profile);
   return messages.NewGeneralStatus(ok, http.StatusOK), 0, nil;
}
