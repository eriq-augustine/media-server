package messages;

import (
   "com/eriq-augustine/mediaserver/model"
)

type CacheStatus struct {
   Success bool
   Progress *model.EncodeProgress
   Queue []model.EncodeRequest
   RecentEncodes []model.CompleteEncode
}

func NewCacheStatus(progress *model.EncodeProgress, queue []model.EncodeRequest, recentEncodes []model.CompleteEncode) *CacheStatus {
   if (progress != nil) {
      safeProgress := progress.Safe();
      progress = &safeProgress;
   }

   safeQueue := make([]model.EncodeRequest, 0, len(queue));
   for _, request := range(queue) {
      safeQueue = append(safeQueue, request.Safe());
   }

   safeRecentEncodes := make([]model.CompleteEncode, 0, len(recentEncodes));
   for _, recentEncode := range(recentEncodes) {
      safeRecentEncodes = append(safeRecentEncodes, recentEncode.Safe());
   }

   return &CacheStatus{true, progress, safeQueue, safeRecentEncodes};
}
