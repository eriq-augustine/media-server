package cache;

// TODO(eriq): The cache/encode system needs a little re-architecting.
// The interactions between the two is a little wonky.

import (
   "os"
   "path/filepath"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

// {extension: requirements}
var fileRequirementsCache *map[string]CacheRequirements;
// {cacheDir: cahceEntry}
var cache map[string]*model.CacheEntry;

type CacheRequirements struct {
   VideoEncode bool
   Subtitles bool
   Poster bool
}

func init() {
   fileRequirementsCache = nil;
   cache = nil;
}

func (requirements CacheRequirements) RequiresCache() bool {
   return requirements.VideoEncode || requirements.Subtitles || requirements.Poster;
}

// The bool return will be true if the cache is ready (or no cache is required) and the file return is valid.
// If the bool return is false, then the cache needs more time to prep.
// Even if the cache is not ready, you should still take the returned file.
// Other components of the cache may have been filled.
func NegotiateCache(file model.File) (model.File, bool) {
   var cacheReady bool = true;
   var filetypes *map[string]CacheRequirements = loadFileRequirements();

   ext := strings.TrimPrefix(filepath.Ext(file.DirEntry.Name), ".");

   requirements, ok := (*filetypes)[ext];

   // If we don't know about this extension or don't need the cache, then return ok.
   if (!ok || !requirements.RequiresCache()) {
      return file, true;
   }

   cacheEntry := getCacheEntry(file);

   if (requirements.Poster) {
      handlePoster(&file, cacheEntry);
   }

   if (requirements.Subtitles) {
      handleSubtitles(&file, cacheEntry);
   }

   if (requirements.VideoEncode) {
      cacheReady = cacheReady && handleEncode(&file, cacheEntry);
   }

   cacheEntry.Save();

   return file, cacheReady;
}

func handleEncode(file *model.File, cacheEntry *model.CacheEntry) bool {
   if (cacheEntry.Encode != nil) {
      cacheLink, err := util.CacheLink(cacheEntry.Encode.EncodePath);
      if (err == nil) {
         file.CacheLink = &cacheLink;
      }

      return true;
   }

   // No cached encode found, request a new one.
   requestEncode(*file, cacheEntry.Dir);

   return false;
}

func handleSubtitles(file *model.File, cacheEntry *model.CacheEntry) {
   subs := cacheEntry.Subtitles;

   if (subs == nil) {
      subs, err := extractSubtitles(file, cacheEntry.Dir);
      if (err != nil) {
         subs = nil;
      } else {
         cacheEntry.SetSubtitles(subs);
      }
   }

   if (subs != nil) {
      for _, sub := range(*subs) {
         subLink, err := util.CacheLink(sub);
         if (err == nil) {
            file.Subtitles = append(file.Subtitles, subLink);
         }
      }
   }
}

func handlePoster(file *model.File, cacheEntry *model.CacheEntry) {
   posterPath := cacheEntry.Poster;

   if (posterPath == nil) {
      posterPath, err := fetchPoster(file, cacheEntry.Dir);
      if (err != nil) {
         posterPath = nil;
      } else {
         cacheEntry.SetPoster(posterPath);
      }
   }

   if (posterPath != nil) {
      posterLink, err := util.CacheLink(*posterPath);
      if (err == nil) {
         file.Poster = &posterLink;
      }
   }
}

func addEncodeToCache(cacheDir string, encode *model.CompleteEncode) {
   if (cache != nil) {
      cache[cacheDir].SetEncode(encode);
      cache[cacheDir].Save();
   }
}

// Note that this causes a race condition between the time the entry is fetched and when the cache is served.
func getCacheEntry(file model.File) *model.CacheEntry {
   if (cache == nil) {
      cache = make(map[string]*model.CacheEntry);
      scanCache();
   }

   cacheDir := ensureCacheDir(file);

   entry, ok := cache[cacheDir];
   if (ok) {
      entry.Hit();
      return entry;
   }

   entry = model.NewCacheEntry(cacheDir);
   cache[cacheDir] = entry;

   return entry;
}

// Scan the cache directory for cache entries and load them into memory.
func scanCache() {
   cachePath := config.GetString("cacheBaseDir");
   err := filepath.Walk(cachePath, func(childPath string, fileInfo os.FileInfo, err error) error {
      if (err != nil) {
         return err;
      }

      if (cachePath == childPath) {
         return nil;
      }

      if (!fileInfo.IsDir() && fileInfo.Name() == model.CACHE_ENTRY_FILE_NAME) {
         cacheEntry := model.LoadCacheEntryFromFile(childPath);
         if (cacheEntry != nil) {
            cache[cacheEntry.Dir] = cacheEntry;
         }
      }

      return nil;
   });

   if (err != nil) {
      log.ErrorE("Error scanning cache", err);
   }
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

func loadFileRequirements() *map[string]CacheRequirements {
   if (fileRequirementsCache != nil) {
      return fileRequirementsCache;
   }

   var requirements map[string]CacheRequirements = make(map[string]CacheRequirements);
   var filetypes map[string]interface{} = config.Get("filetypes").(map[string]interface{});

   for ext, info := range(filetypes) {
      var typeRequirements CacheRequirements;

      cacheInfo, ok := info.(map[string]interface{})["cache"];
      if (ok) {
         cacheContents, ok := cacheInfo.(map[string]interface{})["contents"];
         if (ok) {
            for _, val := range(cacheContents.([]interface{})) {
               switch val.(string) {
               case "video_encode":
                  typeRequirements.VideoEncode = true;
                  break;
               case "subtitles":
                  typeRequirements.Subtitles = true;
                  break;
               case "poster":
                  typeRequirements.Poster = true;
                  break;
               default:
                  log.Error("Unknown cache requirement: " + val.(string));
                  break;
               }
            }
         }
      }

      requirements[ext] = typeRequirements;
   }

   return &requirements;
}
