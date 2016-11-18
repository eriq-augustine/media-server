package cache;

import (
   "os"
   "path/filepath"
   "sort"
   "sync"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

const (
   BYTES_PER_GIGABYTE = 1024 * 1024 * 1024
)

// {cacheDir: cahceEntry}
var cache *map[string]CacheEntry;
var cacheLock sync.Mutex;
var saveLock sync.Mutex;

// This is the worst case possible cache size.
// Instead of checking the cache size every time, maintain this and check for maintenance when
// it grows too large.
// Like |cache|, this should only be accessed under the protection of |cacheLock|.
var maxCacheSize uint64;

func init() {
   cache = nil;
   maxCacheSize = 0;
}

// Load the cache from disk.
func Load() {
   go scanCache();
}

// Use these getters and setters to interact with he cache.

func getCachedPoster(cacheDir string) (string, bool) {
   cacheEntry, ok := internalGetCacheEntry(cacheDir);

   if (!ok || cacheEntry.Poster == nil) {
      return "", false;
   }

   return *cacheEntry.Poster, true;
}

func getCachedSubtitles(cacheDir string) ([]string, bool) {
   cacheEntry, ok := internalGetCacheEntry(cacheDir);

   if (!ok || cacheEntry.Subtitles == nil) {
      return []string{}, false;
   }

   return *cacheEntry.Subtitles, true;
}

func getCachedVideoEncode(cacheDir string) (model.CompleteEncode, bool) {
   cacheEntry, ok := internalGetCacheEntry(cacheDir);

   if (!ok || cacheEntry.VideoEncode == nil) {
      return model.CompleteEncode{}, false;
   }

   return *cacheEntry.VideoEncode, true;
}

func setCachedPoster(cacheDir string, posterPath string) {
   cacheEntry, ok := internalGetCacheEntry(cacheDir);

   if (!ok) {
      cacheEntry = NewCacheEntry(cacheDir);
   }

   cacheEntry.SetPoster(&posterPath);
   internalSetCacheEntry(cacheDir, cacheEntry);
   saveCache(cacheDir);
}

func setCachedSubtitles(cacheDir string, subtitles []string) {
   cacheEntry, ok := internalGetCacheEntry(cacheDir);

   if (!ok) {
      cacheEntry = NewCacheEntry(cacheDir);
   }

   cacheEntry.SetSubtitles(&subtitles);
   internalSetCacheEntry(cacheDir, cacheEntry);
   saveCache(cacheDir);
}

// Update a cache entry after an encode is complete.
func setCachedVideoEncode(cacheDir string, videoEncode model.CompleteEncode) {
   cacheEntry, ok := internalGetCacheEntry(cacheDir);

   if (!ok) {
      cacheEntry = NewCacheEntry(cacheDir);
   }

   cacheEntry.SetVideoEncode(&videoEncode);
   internalSetCacheEntry(cacheDir, cacheEntry);
   saveCache(cacheDir);
}

// Calling this will persist the cache entry to disk.
// The caller should call this after it is fully done with the cache entry.
func saveCache(cacheDir string) {
   saveLock.Lock();
   defer saveLock.Unlock();

   cacheEntry, ok := internalGetCacheEntry(cacheDir);

   if (!ok) {
      return;
   }

   // If the cache is size dirty, then it may have changed sizes recently.
   // Add the cache's size to the max cache size.
   sizeDirty := cacheEntry.IsSizeDirty();

   cacheEntry.Save();

   if (sizeDirty) {
      maxCacheSize += cacheEntry.Size;
   }

   maintainCacheSize();
}

// This is internal only.
// Callers should use specific functions like getCachedPoster() instead.
func internalGetCacheEntry(cacheDir string) (CacheEntry, bool) {
   if (cache == nil) {
      scanCache();
   }

   cacheLock.Lock();
   defer cacheLock.Unlock();

   entry, ok := (*cache)[cacheDir];
   return entry, ok;
}

// Same access rules as internalGetCacheEntry().
func internalSetCacheEntry(cacheDir string, cacheEntry CacheEntry) {
   if (cache == nil) {
      scanCache();
   }

   cacheLock.Lock();
   defer cacheLock.Unlock();

   (*cache)[cacheDir] = cacheEntry;
}

// Scan the cache directory for cache entries and load them into memory.
func scanCache() {
   cacheLock.Lock();
   defer cacheLock.Unlock();

   cache = nil;
   maxCacheSize = 0;

   cacheScan := make(map[string]CacheEntry);
   completeEncodes := make([]model.CompleteEncode, 0);

   cachePath := config.GetString("cacheBaseDir");
   err := filepath.Walk(cachePath, func(childPath string, fileInfo os.FileInfo, err error) error {
      if (err != nil) {
         return err;
      }

      if (cachePath == childPath) {
         return nil;
      }

      if (!fileInfo.IsDir() && fileInfo.Name() == CACHE_ENTRY_FILE_NAME) {
         cacheEntry := LoadCacheEntryFromFile(childPath);
         if (cacheEntry != nil) {
            cacheScan[cacheEntry.Dir] = *cacheEntry;
            maxCacheSize += cacheEntry.Size;

            if (cacheEntry.VideoEncode != nil) {
               completeEncodes = append(completeEncodes, *cacheEntry.VideoEncode);
            }
         }
      }

      return nil;
   });

   if (err != nil) {
      log.ErrorE("Error scanning cache", err);
   }

   cache = &cacheScan;

   loadRecentVideoEncodes(completeEncodes);

   maintainCacheSize();
}

// Note that this does not obtain the cache lock itself.
// The caller should obtain the lock and release it when fully done.
func maintainCacheSize() {
   var totalSize uint64 = 0;

   var lowerThresholdBytes uint64 = uint64(config.GetInt("cacheLowerThresholdGB") * BYTES_PER_GIGABYTE);
   var upperThresholdBytes uint64 = uint64(config.GetInt("cacheUpperThresholdGB") * BYTES_PER_GIGABYTE);

   if (maxCacheSize < upperThresholdBytes) {
      return;
   }

   var sortedCache []CacheEntry = make([]CacheEntry, 0, len(*cache));

   for _, cacheEntry := range(*cache) {
      totalSize += cacheEntry.Size;
      sortedCache = append(sortedCache, cacheEntry);
   }

   // We overestimated too much.
   if (totalSize < upperThresholdBytes) {
      maxCacheSize = totalSize;
      return;
   }

   sort.Sort(CacheByScore(sortedCache));
   for (totalSize > lowerThresholdBytes) {
      if (len(sortedCache) == 0) {
         break;
      }

      cacheEntry := sortedCache[0];
      sortedCache = sortedCache[1:];

      delete(*cache, cacheEntry.Dir);
      removeEncode(cacheEntry.Dir);
      util.RmDir(cacheEntry.Dir);
      totalSize -= cacheEntry.Size;
   }

   maxCacheSize = totalSize;
}
