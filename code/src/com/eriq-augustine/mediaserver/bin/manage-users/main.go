package main;

import (
   "bufio"
   "fmt"
   "os"
   "strings"

   "github.com/howeyc/gopass"
   "golang.org/x/crypto/bcrypt"

   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/user"
   "com/eriq-augustine/mediaserver/util"
);

var reader *bufio.Reader = bufio.NewReader(os.Stdin);

func showUsage() {
   fmt.Println("Manage a users file.\n");
   fmt.Printf("usage: %s <action> <users file>\n\n", os.Args[0]);
   fmt.Println("Options:");
   fmt.Println("   list (ls)   - list the users present the the given file");
   fmt.Println("   add (a)     - add a user to the given file (will create the file if it does not exist)");
   fmt.Println("   remove (rm) - remove a user from the given file");
   fmt.Println("   help (h)    - print this message and exit");
}

func readLine() string {
   text, err := reader.ReadString('\n')
   if (err != nil) {
      fmt.Println("Error reading line: " + err.Error());
      os.Exit(1);
   }
   return strings.TrimSpace(text);
}

func readPassword() string {
   pass, err := gopass.GetPasswd();
   if (err != nil) {
      panic(fmt.Sprintf("Failed to get passowrd: %v", err));
   }

   return strings.TrimSpace(string(pass));
}

func readBool(defaultValue bool) bool {
   stringValue := strings.ToLower(readLine());

   if (stringValue == "") {
      return defaultValue;
   } else if (stringValue == "y" || stringValue == "yes" || stringValue == "t" || stringValue == "true") {
      return true;
   } else if (stringValue == "n" || stringValue == "no" || stringValue == "f" || stringValue == "false") {
      return false;
   } else {
      fmt.Printf("Bad boolean value: %s.\nExiting\n", stringValue);
      os.Exit(1);
      return false;
   }
}

func showListing(usersFile string) {
   if (!util.PathExists(usersFile)) {
      fmt.Printf("Users file (%s) does not exist.\n", usersFile);
      return;
   }

   usersMap := auth.LoadUsersFromFile(usersFile);

   fmt.Printf("User Count: %d\n", len(usersMap));
   for _, user := range(usersMap) {
      fmt.Println("   " + user.String());
   }
}

func addUser(usersFile string) {
   var usersMap map[string]user.User;
   if (util.PathExists(usersFile)) {
      usersMap = auth.LoadUsersFromFile(usersFile);
   } else {
      usersMap = make(map[string]user.User);
   }

   fmt.Print("Username: ");
   username := readLine();

   fmt.Print("Password: ");
   passhash := util.Passhash(username, readPassword());

   fmt.Print("Is Admin [y/N]: ");
   isAdmin := readBool(false);

   bcryptHash, err := bcrypt.GenerateFromPassword([]byte(passhash), bcrypt.DefaultCost);
   if (err != nil) {
      fmt.Printf("Could not generate bcrypt hash: %g", err);
      os.Exit(1);
   }

   usersMap[username] = user.User{username, string(bcryptHash), isAdmin};
   auth.SaveUsersFile(usersFile, usersMap);
}

func removeUser(usersFile string) {
   if (!util.PathExists(usersFile)) {
      fmt.Printf("Users file (%s) does not exist.\n", usersFile);
      return;
   }

   usersMap := auth.LoadUsersFromFile(usersFile);

   fmt.Print("Username: ");
   username := readLine();

   _, exists := usersMap[username];
   if (!exists) {
      fmt.Printf("User (%s) does not exist. Exiting...", username);
      os.Exit(1);
   }

   delete(usersMap, username);
   auth.SaveUsersFile(usersFile, usersMap);
}

func main() {
   args := os.Args;

   if (len(os.Args) != 3 || util.SliceHasString(args, "help") || util.SliceHasString(args, "h")) {
      showUsage();
      return;
   }

   switch args[1] {
   case "list", "ls":
      showListing(args[2]);
      break;
   case "add", "a":
      addUser(args[2]);
      break;
   case "remove", "rm":
      removeUser(args[2]);
      break;
   default:
      fmt.Printf("Unknown action (%s)\n", args[1]);
      showUsage();
   }
}
