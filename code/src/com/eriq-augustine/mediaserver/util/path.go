package util;

import (
   "crypto/rand"
   "fmt"
   "path/filepath"
   "os"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
);

const (
   RANDOM_NAME_LENGTH = 64
   RANDOM_CHARS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Get the basename for |path|.
// That is, the name of the file (last component) without any extension.
func Basename(path string) string {
   ext := filepath.Ext(path);
   return strings.TrimSuffix(filepath.Base(path), ext);
}

// Tell if a path exists.
func PathExists(path string) bool {
   _, err := os.Stat(path);
   if (err != nil) {
      if os.IsNotExist(err) {
         return false;
      }
   }

   return true;
}

// You should probably pass an extension, but the other two can be empty strings.
func TempFilePath(extension string, prefix string, suffix string) string {
   filename := prefix + RandomString(RANDOM_NAME_LENGTH) + suffix + "." + extension;
   return filepath.Join(os.TempDir(), filename);
}

func RandomString(length int) string {
   bytes := make([]byte, length);
   _, err := rand.Read(bytes);
   if (err != nil) {
      log.ErrorE("Unable to generate random string", err);
   }

   for i, val := range(bytes) {
      bytes[i] = RANDOM_CHARS[int(val) % len(RANDOM_CHARS)];
   }

   return string(bytes)
}

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
func AbstractPath(realPath string, baseDir string) (string, error) {
   if (baseDir == "") {
      baseDir = config.GetString("staticBaseDir");
   }

   abstractPath, err := filepath.Rel(baseDir, realPath);
   if (err != nil) {
      return "", err;
   }

   return abstractPath, nil;
}

// Given a clean path, get the link to that resource.
func RawLink(path string) (string, error) {
   abstractPath, err := AbstractPath(path, config.GetString("staticBaseDir"));
   if (err != nil) {
      return "", err;
   }

   return filepath.Join("/", config.GetString("rawBaseURL"), abstractPath), nil;
}

func CacheLink(path string) (string, error) {
   abstractPath, err := AbstractPath(path, config.GetString("cacheBaseDir"));
   if (err != nil) {
      return "", err;
   }

   return filepath.Join("/", config.GetString("cacheBaseURL"), abstractPath), nil;
}
