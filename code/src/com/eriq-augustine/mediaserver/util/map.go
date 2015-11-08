package util;

// Some utilities for maps.

func MapHasKey(haystack map[string]string, needle string) bool {
   _, ok := haystack[needle];
   return ok;
}
