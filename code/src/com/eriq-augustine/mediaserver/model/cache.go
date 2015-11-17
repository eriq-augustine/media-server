package model;

import (
   "encoding/json"
   "io/ioutil"
   "path/filepath"
   "time"

   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/util"
)

const (
   CACHE_ENTRY_FILE_NAME = "cache.json"
)

type CacheEntry struct {
   Dir string
   Size uint64
   Hits int
   LastHit time.Time
   LastUpdate time.Time
   Subtitles *[]string
   Poster *string
   Encode *CompleteEncode

   dirty bool // If the entry is dirty, it needs to be saved.
   sizeDirty bool // Whether or not the size needs to be recalculated.
}

func NewCacheEntry(dir string) *CacheEntry {
   return &CacheEntry {
      Dir: dir,
      Size: 0,
      Hits: 1,
      LastHit: time.Now(),
      LastUpdate: time.Now(),
      Subtitles: nil,
      Poster: nil,
      Encode: nil,
      dirty: true,
   };
}

func LoadCacheEntryFromFile(path string) *CacheEntry {
   data, err := ioutil.ReadFile(path);
   if (err != nil) {
      log.ErrorE("Unable to read cache entry", err);
      return nil;
   }

   var cacheEntry CacheEntry;
   err = json.Unmarshal(data, &cacheEntry);
   if (err != nil) {
      log.ErrorE("Unable to unmarshal cache entry", err);
      return nil;
   }

   return &cacheEntry;
}

func (entry *CacheEntry) Hit() {
   entry.Hits++;
   entry.LastHit = time.Now();
   entry.dirty = true;
}

func (entry *CacheEntry) SetPoster(poster *string) {
   entry.Poster = poster;
   entry.LastUpdate = time.Now();
   entry.dirty = true;
   entry.sizeDirty = true;
}

func (entry *CacheEntry) SetSubtitles(subs *[]string) {
   entry.Subtitles = subs;
   entry.LastUpdate = time.Now();
   entry.dirty = true;
   entry.sizeDirty = true;
}

func (entry *CacheEntry) SetEncode(encode *CompleteEncode) {
   entry.Encode = encode;
   entry.LastUpdate = time.Now();
   entry.dirty = true;
   entry.sizeDirty = true;
}

func (entry *CacheEntry) Save() {
   if (!entry.dirty) {
      return;
   }

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

type EncodeRequest struct {
   File File
   CacheDir string
   RequestTime time.Time
}

// Get a "safe" copy (one with no sensitive information).
func (request EncodeRequest) Safe() EncodeRequest {
   return EncodeRequest{
      File: request.File.Safe(),
      CacheDir: "",
      RequestTime: request.RequestTime,
   };
}

type EncodeProgress struct {
   File File
   CacheDir string
   CompleteMS int64
   TotalMS int64
   Done bool
}

func (progress EncodeProgress) Safe() EncodeProgress {
   return EncodeProgress {
      File: progress.File.Safe(),
      CacheDir: "",
      CompleteMS: progress.CompleteMS,
      TotalMS: progress.TotalMS,
      Done: progress.Done,
   };
}

type CompleteEncode struct {
   File File
   CompleteTime time.Time
   EncodePath string
}

func (complete CompleteEncode) Safe() CompleteEncode {
   return CompleteEncode{
      File: complete.File.Safe(),
      CompleteTime: complete.CompleteTime,
      EncodePath: "",
   };
}
