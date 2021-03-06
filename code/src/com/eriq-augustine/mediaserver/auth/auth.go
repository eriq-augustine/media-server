package auth;

// We use fairly simple security.
// Auth the user and give them a token that does not expire.
// However the token is stored in memory, so a server restart invalidates it.

import (
   "bytes"
   "crypto/rand"
   "encoding/base64"
   "encoding/binary"
   "encoding/json"
   "fmt"
   "io/ioutil"
   "sync"
   "time"

   "github.com/eriq-augustine/elfs-api/apierrors"
   "github.com/eriq-augustine/goconfig"
   "github.com/eriq-augustine/golog"
   "golang.org/x/crypto/bcrypt"

   "com/eriq-augustine/mediaserver/user"
   "com/eriq-augustine/mediaserver/util"
)

const (
   TOKEN_RANDOM_BYTE_LEN = 16
)

// {username: User}
var Users map[string]user.User
// {token: username}
var Sessions map[string]string;

var createAccountMutex sync.Mutex;

func init() {
   Users = make(map[string]user.User);
   Sessions = make(map[string]string);
}

// Returns the token.
func AuthenticateUser(username string, passhash string) (string, error) {
   _, exists := Users[username];
   if (!exists) {
      return "", apierrors.TokenValidationError{apierrors.TOKEN_AUTH_BAD_CREDENTIALS};
   }

   err := bcrypt.CompareHashAndPassword([]byte(Users[username].Passhash), []byte(passhash));
   if (err != nil) {
      return "", apierrors.TokenValidationError{apierrors.TOKEN_AUTH_BAD_CREDENTIALS};
   }

   token, _:= generateToken();
   Sessions[token] = username;

   return token, nil;
}

// We just keep tokens in memory.
func RegisterToken(username string, token string) error {
   return nil;
}

// Validate the token and get back the token's secret.
func ValidateToken(token string) (string, error) {
   username, exists := Sessions[token];
   if (!exists) {
      return "", apierrors.TokenValidationError{apierrors.TOKEN_VALIDATION_NO_TOKEN};
   }

   return username, nil;
}

// Invalidate the token.
func InvalidateToken(token string) (bool, error) {
   _, exists := Sessions[token];
   if (!exists) {
      return false, apierrors.TokenValidationError{apierrors.TOKEN_VALIDATION_NO_TOKEN};
   }

   delete(Sessions, token);
   return true, nil;
}

func CreateUser(username string, passhash string) (string, error) {
   bcryptHash, err := bcrypt.GenerateFromPassword([]byte(passhash), bcrypt.DefaultCost);
   if (err != nil) {
      golog.ErrorE("Could not generate bcrypt hash", err);
      return "", err;
   }

   createAccountMutex.Lock();
   defer createAccountMutex.Unlock();

   _, exists := Users[username];
   if (exists) {
      return "", fmt.Errorf("Username (%s) already exists", username);
   }

   Users[username] = user.User{username, string(bcryptHash), false};

   token, _:= generateToken();
   Sessions[token] = username;

   SaveUsers();

   return token, nil;
}

func SaveUsers() {
   SaveUsersFile(goconfig.GetString("usersFile"), Users);
}

func SaveUsersFile(usersFile string, usersMap map[string]user.User) {
   fileUsers := make([]user.User, 0);
   for _, user := range(usersMap) {
      fileUsers = append(fileUsers, user);
   }

   jsonString, err := util.ToJSONPretty(fileUsers);
   if (err != nil) {
      golog.ErrorE("Unable to marshal users", err);
      return;
   }

   err = ioutil.WriteFile(usersFile, []byte(jsonString), 0644);
   if (err != nil) {
      golog.ErrorE("Unable to save users", err);
   }
}

func LoadUsers() {
   Users = LoadUsersFromFile(goconfig.GetString("usersFile"));
}

func LoadUsersFromFile(usersFile string) map[string]user.User {
   usersMap := make(map[string]user.User);

   data, err := ioutil.ReadFile(usersFile);
   if (err != nil) {
      golog.ErrorE("Unable to read users file", err);
      return usersMap;
   }

   var fileUsers []user.User;
   err = json.Unmarshal(data, &fileUsers);
   if (err != nil) {
      golog.ErrorE("Unable to unmarshal users", err);
      return usersMap;
   }

   for _, user := range(fileUsers) {
      usersMap[user.Username] = user;
   }

   return usersMap;
}

// Generate a "unique" token.
// Return the base64 encoding of the token as well as the time it was created.
// The token is a base64 encoding of <date (micro unix epoch), user id, rand>.
// The date takes 8 bytes (uint64) and the random section is TOKEN_RANDOM_BYTE_LEN bytes.
func generateToken() (string, time.Time) {
   now := time.Now();

   randData := make([]byte, TOKEN_RANDOM_BYTE_LEN);
   rand.Read(randData);

   timeBinary := make([]byte, 8);
   binary.LittleEndian.PutUint64(timeBinary, uint64(now.UnixNano() / 1000));

   tokenData := bytes.Join([][]byte{timeBinary, randData}, []byte{});

   return base64.URLEncoding.EncodeToString(tokenData), now;
}
