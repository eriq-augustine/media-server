package model;

import (
   "time"
)

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
}

func (complete CompleteEncode) Safe() CompleteEncode {
   return CompleteEncode{
      File: complete.File.Safe(),
      CompleteTime: complete.CompleteTime,
   };
}
