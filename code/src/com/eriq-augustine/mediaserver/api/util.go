package api;

import (
   "net/http"
   "strings"
);

func getToken(request *http.Request) (string, bool) {
   authHeader, ok := request.Header["Authorization"];

   if (!ok) {
      return "", false;
   }

   if (len(authHeader) == 0) {
      return "", false;
   }

   token := strings.TrimPrefix(strings.TrimSpace(authHeader[0]), "Bearer ");

   if (token == "") {
      return "", false;
   }

   return token, true;
}
