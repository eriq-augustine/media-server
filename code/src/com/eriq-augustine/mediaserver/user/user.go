package user;

import (
   "fmt"
)

type User struct {
   Username string
   Passhash string
   IsAdmin bool
}

func (user User) String() string {
   if (user.IsAdmin) {
      return fmt.Sprintf("%s (Admin)", user.Username);
   }

   return user.Username;
}
