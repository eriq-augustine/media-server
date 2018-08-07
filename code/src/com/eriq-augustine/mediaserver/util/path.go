package util;

import (
   "crypto/rand"
   "fmt"
   "path/filepath"
   "os"
   "strings"

   "github.com/eriq-augustine/elfs/dirent"
   "github.com/eriq-augustine/elfs-api/model"
   "github.com/eriq-augustine/goconfig"
   "github.com/eriq-augustine/golog"
);

const (
   RANDOM_NAME_LENGTH = 64
   RANDOM_CHARS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func DirEntryFromInfo(fileInfo os.FileInfo, path string) *model.DirEntry {
   abstractPath, err := AbstractPath(path);
   if (err != nil) {
      golog.WarnE("Could not get abstract path for (" + path + ")", err);
      abstractPath = path;
   }

   abstractParentPath, err := AbstractPath(ParentPath(path));
   if (err != nil) {
      golog.WarnE("Could not get abstract path for parent.", err);
      abstractParentPath = ParentPath(path);
   }

   var name string = BasenameWithExt(path);
   if (abstractPath == "") {
      name = "";
   }

   return &model.DirEntry{
      Id: dirent.Id(abstractPath),
      IsFile: !fileInfo.IsDir(),
      Owner: 0,
      Name: name,
      CreateTimestamp: fileInfo.ModTime().Unix(),
      ModTimestamp: fileInfo.ModTime().Unix(),
      AccessTimestamp: fileInfo.ModTime().Unix(),
      AccessCount: 1,
      GroupPermissions: nil,
      Size: uint64(fileInfo.Size()),
      Md5: "",
      Parent: dirent.Id(abstractParentPath),
   };
}

// Get the basename for |path|.
// That is, the name of the file (last component) without any extension.
func Basename(path string) string {
   ext := filepath.Ext(path);
   return strings.TrimSuffix(filepath.Base(path), ext);
}

func BasenameWithExt(path string) string {
   return filepath.Base(path);
}

func Ext(path string) string {
   return strings.TrimPrefix(filepath.Ext(path), ".");
}

func ParentPath(path string) string {
   return filepath.Dir(filepath.Clean(path));
}

func Join(parts... string) string {
   return filepath.Clean(filepath.Join(parts...));
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
      golog.ErrorE("Unable to generate random string", err);
   }

   for i, val := range(bytes) {
      bytes[i] = RANDOM_CHARS[int(val) % len(RANDOM_CHARS)];
   }

   return string(bytes)
}

// Take in an abstract path from the clinet and convert it into a real path.
func RealPath(path string) (string, error) {
   cleanPath := filepath.Join(goconfig.GetString("staticBaseDir"), strings.TrimPrefix(path, "/"));

   cleanPath, err := filepath.Abs(cleanPath);
   if (err != nil) {
      return "", err;
   }

   cleanPath = filepath.Clean(cleanPath);

   // Ensure that the path is inside the root directory.
   relPath, err := filepath.Rel(goconfig.GetString("staticBaseDir"), cleanPath);
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
   var baseDir string = goconfig.GetString("staticBaseDir");

   abstractPath, err := filepath.Rel(baseDir, realPath);
   if (err != nil) {
      return "", err;
   }

   // Correct for root and outside paths.
   if (abstractPath == "." || strings.HasPrefix(abstractPath, "..")) {
      abstractPath = "";
   }

   return abstractPath, nil;
}
