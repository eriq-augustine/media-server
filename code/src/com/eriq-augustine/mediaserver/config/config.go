package config;

// For the defaulted getters, the defualt will be returned on ANY error (even if the key exists, but is of the
// wrong type.

import (
   "encoding/json"
   "fmt"
   "math"
   "os"

   "github.com/eriq-augustine/golog"
);

var options map[string]interface{} = make(map[string]interface{});

// Load a file into configuration.
// This will not clear out an existing configuration (so can load multiple files).
// If there are any key conflicts, the file loaded last will win.
// If you want to clear the config, use Reset().
func LoadFile(filename string) (bool, error) {
   file, err := os.Open(filename);

   if (err != nil) {
      golog.ErrorE("Could not open config file", err);
      return false, err;
   }
   defer file.Close();

   decoder := json.NewDecoder(file);

   var fileOptions map[string]interface{};

   err = decoder.Decode(&fileOptions);
   if (err != nil) {
      golog.ErrorE("Unable to decode config file: " + filename, err);
      return false, err;
   }

   for key, val := range fileOptions {
      // encoding/json uses float64 as its default numeric type.
      // Check if it is actually an integer.
      floatVal, ok := val.(float64);
      if (ok) {
         if (math.Trunc(floatVal) == floatVal) {
            val = int(floatVal);
         }
      }
      options[key] = val;
   }

   return true, nil;
}

func Reset() {
   options = make(map[string]interface{});
}

func Has(key string) bool {
   _, present := options[key];
   return present;
}

func Get(key string) interface{} {
   val, present := options[key];
   if (!present) {
      golog.Panic(fmt.Sprintf("Option (%s) does not exist", key));
   }

   return val;
}

func GetDefault(key string, defaultVal interface{}) interface{} {
   if (!Has(key)) {
      return defaultVal;
   }

   val, _ := options[key];
   return val;
}

func GetString(key string) string {
   val := Get(key);

   stringVal, ok := val.(string);
   if (!ok) {
      golog.Panic(fmt.Sprintf("Option (%s) is not a string type", key));
   }

   return stringVal;
}

func GetStringDefault(key string, defaultVal string) string {
   val := GetDefault(key, defaultVal);

   stringval, ok := val.(string);
   if (!ok) {
      return defaultVal;
   }

   return stringval;
}

func GetInt(key string) int {
   val := Get(key);

   intVal, ok := val.(int);
   if (!ok) {
      golog.Panic(fmt.Sprintf("Option (%s) is not an int type", key));
   }

   return intVal;
}

func GetIntDefault(key string, defaultVal int) int {
   val := GetDefault(key, defaultVal);

   intVal, ok := val.(int);
   if (!ok) {
      return defaultVal;
   }

   return intVal;
}

func GetBool(key string) bool {
   val := Get(key);

   boolVal, ok := val.(bool);
   if (!ok) {
      golog.Panic(fmt.Sprintf("Option (%s) is not a bool type", key));
   }

   return boolVal;
}

func GetBoolDefault(key string, defaultVal bool) bool {
   val := GetDefault(key, defaultVal);

   boolVal, ok := val.(bool);
   if (!ok) {
      return defaultVal;
   }

   return boolVal;
}
