package cache;

import (
   "os"
   "path/filepath"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

var fileRequirementsCache *map[string]CacheRequirements = nil;

// TODO(eriq): Right now, all video encodes will be mp4.
// We may change this later.
type CacheRequirements struct {
   VideoEncode bool
   Subtitles bool
   Poster bool
}

func (requirements CacheRequirements) RequiresCache() bool {
   return requirements.VideoEncode || requirements.Subtitles || requirements.Poster;
}

func NegotiateCache(file model.File) {
   var filetypes *map[string]CacheRequirements = loadFileRequirements();

   ext := strings.TrimPrefix(filepath.Ext(file.DirEntry.Name), ".");

   requirements, ok := (*filetypes)[ext];
   if (ok && requirements.RequiresCache()) {
      cacheDir := ensureCacheDir(file);

      if (requirements.Poster) {
         fetchPoster(file, cacheDir);
      }

      if (requirements.Subtitles) {
         extractSubtitles(file, cacheDir);
      }

      if (requirements.VideoEncode) {
         requestEncode(file, cacheDir);
      }
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
