package database;

// Handle the data store related operations.
// The implementation of this does not necessarialy need to be backed with a database.

import (
   "time"

   "com/eriq-augustine/mediaserver/log"
)

// Returns the user's id.
func AuthenticateUser(username string, passhash string) (int, error) {
   log.Debug("Auth");
   return 0, nil;
}

func RegisterToken(userId int, secret string, tokenString string, createTime time.Time, expireTime time.Time) error {
   return nil;
}

// Validate the token and get back the token's secret.
func ValidateToken(tokenString string, userId int) (string, error) {
   return "", nil;
}

// Invalidate the token.
func InvalidateToken(tokenString string) (bool, error) {
   return false, nil;
}
