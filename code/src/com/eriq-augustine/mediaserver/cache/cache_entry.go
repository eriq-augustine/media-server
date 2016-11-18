package cache;

import (
   "encoding/json"
   "fmt"
   "io/ioutil"
   "path/filepath"
   "time"

   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

const (
   CACHE_ENTRY_FILE_NAME = "cache.json"
)

type CacheByScore []CacheEntry

type CacheEntry struct {
   Dir string
   Size uint64
   Hits int // Hits do not actually reference direct cache hits. Instead, it is the number of saves.
   LastHit time.Time
   LastUpdate time.Time
   Subtitles *[]string
   Poster *string
   VideoEncode *model.CompleteEncode

   dirty bool // If the entry is dirty, it needs to be saved.
   sizeDirty bool // Whether or not the size needs to be recalculated.
   score float64
}

func NewCacheEntry(dir string) CacheEntry {
   return CacheEntry {
      Dir: dir,
      Size: 0,
      Hits: 0,
      LastHit: time.Now(),
      LastUpdate: time.Now(),
      Subtitles: nil,
      Poster: nil,
      VideoEncode: nil,
      dirty: true,
      score: 0,
   };
}

func LoadCacheEntryFromFile(path string) *CacheEntry {
   data, err := ioutil.ReadFile(path);
   if (err != nil) {
      log.ErrorE(fmt.Sprintf("Unable to read cache entry: %s", path), err);
      return nil;
   }

   var cacheEntry CacheEntry;
   err = json.Unmarshal(data, &cacheEntry);
   if (err != nil) {
      log.ErrorE(fmt.Sprintf("Unable to unmarshal cache entry: %s", path), err);
      return nil;
   }

   return &cacheEntry;
}

func (entry *CacheEntry) SetPoster(poster *string) {
   entry.Poster = poster;
   entry.LastUpdate = time.Now();
   entry.dirty = true;
   entry.sizeDirty = true;
   entry.score = 0;
}

func (entry *CacheEntry) SetSubtitles(subs *[]string) {
   entry.Subtitles = subs;
   entry.LastUpdate = time.Now();
   entry.dirty = true;
   entry.sizeDirty = true;
   entry.score = 0;
}

func (entry *CacheEntry) SetVideoEncode(encode *model.CompleteEncode) {
   entry.VideoEncode = encode;
   entry.LastUpdate = time.Now();
   entry.dirty = true;
   entry.sizeDirty = true;
   entry.score = 0;
}

func (entry *CacheEntry) IsSizeDirty() bool {
   return entry.sizeDirty;
}

func (entry *CacheEntry) Save() {
   if (!entry.dirty) {
      return;
   }

   entry.Hits++;

   if (entry.sizeDirty) {
      entry.Size = util.DirSize(entry.Dir);
      entry.sizeDirty = false;
   }

   jsonString, err := util.ToJSONPretty(entry);
   if (err != nil) {
      log.ErrorE("Unable to marshal cache entry", err);
      return;
   }

   err = ioutil.WriteFile(filepath.Join(entry.Dir, CACHE_ENTRY_FILE_NAME), []byte(jsonString), 0644);
   if (err != nil) {
      log.ErrorE("Unable to save cache entry", err);
   }

   entry.dirty = false;
}

// High scores are bad.
// The higher the socre, the more likely an entry will be replaced.
func (entry *CacheEntry) Score() float64 {
   if (entry.score != 0) {
      return entry.score;
   }

   dayDiff := uint64(time.Now().Sub(entry.LastHit).Hours() / 24);

   // This is just some made-up scoring formula.
   entry.score = float64(((dayDiff + 1.0) * entry.Size) / ((uint64(entry.Hits) + 1.0) / 2.0));
   return entry.score;
}

func (entries CacheByScore) Len() int {
   return len(entries);
}

func (entries CacheByScore) Swap(i int, j int) {
   entries[i], entries[j] = entries[j], entries[i];
}

func (entries CacheByScore) Less(i int, j int) bool {
   return entries[i].Score() > entries[j].Score()
}
