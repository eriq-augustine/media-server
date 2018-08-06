package main;

import (
   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/server"
);

func main() {
   server.LoadConfig();

   // It is safe to load users after the configs have been loaded.
   auth.LoadUsers();

   server.StartServer();
}
