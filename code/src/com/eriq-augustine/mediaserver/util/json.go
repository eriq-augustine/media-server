package util;

import (
   "encoding/json"

   "com/eriq-augustine/mediaserver/log"
);

func ToJSON(data interface{}) (string, error) {
   bytes, err := json.Marshal(data);
   if (err != nil) {
      log.ErrorE("Error converting to JSON", err);
      return "", err;
   }

   return string(bytes), nil;
}

func ToJSONPretty(data interface{}) (string, error) {
   bytes, err := json.MarshalIndent(data, "", "   ");
   if (err != nil) {
      log.ErrorE("Error converting to JSON", err);
      return "", err;
   }

   return string(bytes), nil;
}
