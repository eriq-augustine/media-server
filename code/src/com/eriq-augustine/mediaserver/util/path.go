package util;

import (
   "fmt"
   "path/filepath"
   "strings"

   "com/eriq-augustine/mediaserver/config"
);

// Take in an abstract path from the clinet and convert it into a real path.
func RealPath(path string) (string, error) {
   cleanPath := filepath.Join(config.GetString("staticBaseDir"), strings.TrimPrefix(path, "/"));

   cleanPath, err := filepath.Abs(cleanPath);
   if (err != nil) {
      return "", err;
   }

   cleanPath = filepath.Clean(cleanPath);

   // Ensure that the path is inside the root directory.
   relPath, err := filepath.Rel(config.GetString("staticBaseDir"), cleanPath);
   if (err != nil) {
      return "", err;
   }

   if (strings.HasPrefix(relPath, "..")) {
      return "", fmt.Errorf("Path outside of root");
   }

   return cleanPath, nil;
}

// Given a real path, find an abstract path for it.
// Abstract paths are just relative paths from the static base directory.
func AbstractPath(realPath string) (string, error) {
   abstractPath, err := filepath.Rel(config.GetString("staticBaseDir"), realPath);
   if (err != nil) {
      return "", err;
   }

   // TEST
   fmt.Println(realPath);
   fmt.Println(abstractPath);

   return abstractPath, nil;
}

// Given a clean path, get the link to that resource.
func RawLink(path string) (string, error) {
   abstractPath, err := AbstractPath(path);
   if (err != nil) {
      return "", err;
   }

   return filepath.Join("/", config.GetString("rawBaseURL"), abstractPath), nil;
}

func CacheLink(path string) (bool, string) {
   // TODO(eriq): Do Caching
   return false, "";

   /*
   abstractPath, err := AbstractPath(path);
   if (err != nil) {
      return "", err;
   }

   return filepath.Join("/", config.GetString("cacheBaseURL"), abstractPath), nil;
   */
}
