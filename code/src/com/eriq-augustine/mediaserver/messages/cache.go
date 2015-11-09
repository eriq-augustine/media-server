package messages;

import (
   "com/eriq-augustine/mediaserver/model"
)

type CacheStatus struct {
   Success bool
   Progress *model.EncodeProgress
   Queue []model.EncodeRequest
   RecentEncodes []model.EncodeRequest
}

func NewCacheStatus(progress *model.EncodeProgress, queue []model.EncodeRequest, recentEncodes []model.EncodeRequest) *CacheStatus {
   return &CacheStatus{true, progress, queue, recentEncodes};
}
