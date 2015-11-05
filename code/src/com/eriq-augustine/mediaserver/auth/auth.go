package auth;

import (
   "crypto/rand"
   "encoding/hex"
   "fmt"
   "time"

   jwt "github.com/dgrijalva/jwt-go"

   "com/eriq-augustine/mediaserver/database"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/util/errors"
);

const (
   SECRET_LENGTH = 256
   TOKEN_EXPIRE_TIME_MIN = 30 * 24 * 60 // One month
)

var jwtSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256;

func genRandomSecret() string {
   secret := make([]byte, SECRET_LENGTH);
   _, err := rand.Read(secret);
   if (err != nil) {
      log.ErrorE("Unable to generate random secret", err);
   }
   return hex.EncodeToString(secret);
}

// Returns a signed jwt string.
func AuthenticateUser(username string, passhash string) (int, string, error) {
   userId, err := database.AuthenticateUser(username, passhash);

   if (err != nil) {
      return -1, "", err;
   }

   secret := genRandomSecret();

   now := time.Now();
   expireTime := now.Add(time.Duration(time.Minute) * TOKEN_EXPIRE_TIME_MIN);

   token := jwt.New(jwtSigningMethod)
   token.Claims["userId"] = userId;
   token.Claims["username"] = username;
   token.Claims["tokenCreateTime"] = now.Unix();
   // Token expire time.
   token.Claims["exp"] = expireTime.Unix();
   tokenString, err := token.SignedString([]byte(secret));

   if (err != nil) {
      log.ErrorE("Unable to sign token", err);
      return -1, "", err;
   }

   err = database.RegisterToken(userId, secret, tokenString, now, expireTime);

   if (err != nil) {
      log.ErrorE("Unable to register token", err);
      return -1, "", err;
   }

   return userId, tokenString, nil;
}

// Return the user id associated with this token on success.
func ValidateToken(tokenString string) (int, error) {
   // The key function will not return our errors back, so we will stach it.
   var tokenError error = nil;
   var userId int = -1;

   // The callback is supposed to return the token's secret.
   // We will validate the token from the database prespective at the same time.
   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      // First validate the signing algorithm.
      if (token.Method.Alg() != jwtSigningMethod.Alg()) {
         return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"]);
      }

      userIdFloat, ok := token.Claims["userId"].(float64);
      if (!ok) {
         // Token has been tampered with.
         tokenError = errors.TokenValidationError{errors.TOKEN_VALIDATION_BAD_SIGNATURE};
         return nil, tokenError;
      }

      userId = int(userIdFloat);
      secret, err := database.ValidateToken(tokenString, userId);
      if (err != nil) {
         tokenError, _ = err.(errors.TokenValidationError);
         return nil, err;
      }

      return []byte(secret), nil;
   });

   if (err != nil) {
      // First check to see if there was a token specific error,
      // then just return the default error.
      if (tokenError != nil) {
         return -1, tokenError;
      }
      return -1, err;
   }

   // Check for expired tokens
   expFloat, ok := token.Claims["exp"].(float64);
   if (!ok) {
      // Token has been tampered with.
      return -1, errors.TokenValidationError{errors.TOKEN_VALIDATION_BAD_SIGNATURE};
   }

   if (time.Unix(int64(expFloat), 0).Before(time.Now())) {
      return -1, errors.TokenValidationError{errors.TOKEN_VALIDATION_EXPIRED};
   }

   if (!token.Valid) {
      return -1, errors.TokenValidationError{errors.TOKEN_VALIDATION_BAD_SIGNATURE};
   }

   return userId, nil;
}

func InvalidateToken(tokenString string) (bool, error) {
   return database.InvalidateToken(tokenString);
}
