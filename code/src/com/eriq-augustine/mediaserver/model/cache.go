package model;

type EncodeRequest struct {
   File File
   CacheDir string
}

type EncodeProgress struct {
   File File
   CacheDir string
   CompleteMS int64
   TotalMS int64
   Done bool
}
