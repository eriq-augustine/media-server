package util;

// Some utilities for maps.

func MapHasKey(haystack map[string]string, needle string) bool {
   _, ok := haystack[needle];
   return ok;
}

func MapGetWithDefault(haystack map[string]string, needle string, defaultValue string) string {
   val, ok := haystack[needle];

   if (ok) {
      return val;
   } else {
      return defaultValue;
   }
}
