package cache;

// The entry point into the caching system.

import (
   "os"
   "path/filepath"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

// The bool return will be true if the cache is ready (or no cache is required) and the file return is valid.
// If the bool return is false, then the cache needs more time to prep.
// Even if the cache is not ready, you should still take the returned file.
// Other components of the cache may have been filled.
func NegotiateCache(file model.File) (model.File, bool) {
   if (!requiresCache(file)) {
      return file, true;
   }

   var cacheReady bool = true;
   var cacheDir string = ensureCacheDir(file);

   if (requiresPoster(file)) {
      handlePoster(&file, cacheDir)
   }

   if (requiresSubtitles(file)) {
      handleSubtitles(&file, cacheDir);
   }

   if (requiresVideoEncode(file)) {
      cacheReady = cacheReady && handleVideoEncode(&file, cacheDir);
   }

   saveCache(cacheDir);

   return file, cacheReady;
}

func GetEncodeStatus() model.EncodeStatus {
   return getEncodeStatus();
}

func handleVideoEncode(file *model.File, cacheDir string) bool {
   encode, ok := getCachedVideoEncode(cacheDir);

   if (ok) {
      encodeLink, err := util.CacheLink(encode.EncodePath);
      if (err == nil) {
         file.CacheLink = &encodeLink;
      }

      return true;
   }

   // No cached encode found, request a new one.
   queueEncode(*file, cacheDir);

   return false;
}

func handleSubtitles(file *model.File, cacheDir string) {
   subs, ok := getCachedSubtitles(cacheDir);

   if (!ok) {
      subs, err := extractSubtitles(file, cacheDir);
      if (err != nil) {
         return;
      }

      setCachedSubtitles(cacheDir, subs);
   }

   subLinks := make([]string, 0, len(subs));
   for _, sub := range(subs) {
      subLink, err := util.CacheLink(sub);
      if (err == nil) {
         subLinks = append(subLinks, subLink);
      }
   }

   file.Subtitles = subLinks;
}

func handlePoster(file *model.File, cacheDir string) {
   posterPath, ok := getCachedPoster(cacheDir);

   if (!ok) {
      posterPath, err := fetchPoster(file, cacheDir);
      if (err != nil) {
         return;
      }

      setCachedPoster(cacheDir, posterPath);
   }

   posterLink, err := util.CacheLink(posterPath);
   if (err != nil) {
      return;
   }

   file.Poster = &posterLink;
}

func ensureCacheDir(file model.File) string {
   cacheName := util.SHA1Hex(file.DirEntry.Path);
   cachePath := filepath.Join(config.GetString("cacheBaseDir"), cacheName)

   err := os.MkdirAll(cachePath, 0755);
   if (err != nil) {
      log.ErrorE("Error ensuring cache dir: " + cachePath, err);
   }

   return cachePath;
}
