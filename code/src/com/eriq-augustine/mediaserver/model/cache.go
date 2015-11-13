package model;

import (
   "time"
)

type EncodeRequest struct {
   File File
   CacheDir string
   RequestTime time.Time
}

type EncodeProgress struct {
   File File
   CacheDir string
   CompleteMS int64
   TotalMS int64
   Done bool
}

type CompleteEncode struct {
   File File
   CompleteTime time.Time
}
