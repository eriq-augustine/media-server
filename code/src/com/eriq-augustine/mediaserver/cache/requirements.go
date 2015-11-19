package cache;

// Handle describing the caching requirements for each filetype.

import (
   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

type CacheRequirements struct {
   VideoEncode bool
   Subtitles bool
   Poster bool
}

func (requirements CacheRequirements) RequiresCache() bool {
   return requirements.VideoEncode || requirements.Subtitles || requirements.Poster;
}

// {extension: requirements}
var fileRequirementsCache *map[string]CacheRequirements;

func init() {
   fileRequirementsCache = nil;
}

func requiresCache(file model.File) bool {
   requirements, ok := getRequirements(file);
   if (ok) {
      return requirements.RequiresCache();
   }

   return false;
}

func requiresPoster(file model.File) bool {
   requirements, ok := getRequirements(file);
   if (ok) {
      return requirements.Poster;
   }

   return false;
}

func requiresSubtitles(file model.File) bool {
   requirements, ok := getRequirements(file);
   if (ok) {
      return requirements.Subtitles;
   }

   return false;
}

func requiresVideoEncode(file model.File) bool {
   requirements, ok := getRequirements(file);
   if (ok) {
      return requirements.VideoEncode;
   }

   return false;
}

func getRequirements(file model.File) (CacheRequirements, bool) {
   if (fileRequirementsCache == nil) {
      loadFileRequirements();
   }

   requirements, ok := (*fileRequirementsCache)[util.Ext(file.DirEntry.Name)];
   if (ok) {
      return requirements, true;
   }

   return CacheRequirements{}, false;
}

func loadFileRequirements() {
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

   fileRequirementsCache = &requirements;
}
