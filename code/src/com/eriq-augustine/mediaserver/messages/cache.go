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

func NewCacheStatus(encodeStatus model.EncodeStatus) *CacheStatus {
   var safeProgress *model.EncodeProgress = nil;
   if (encodeStatus.Progress != nil) {
      temp := encodeStatus.Progress.Safe();
      safeProgress = &temp;
   }

   safeQueue := make([]model.EncodeRequest, 0, len(encodeStatus.Queue));
   for _, request := range(encodeStatus.Queue) {
      safeQueue = append(safeQueue, request.Safe());
   }

   safeRecentEncodes := make([]model.CompleteEncode, 0, len(encodeStatus.Complete));
   for _, recentEncode := range(encodeStatus.Complete) {
      safeRecentEncodes = append(safeRecentEncodes, recentEncode.Safe());
   }

   return &CacheStatus{true, safeProgress, safeQueue, safeRecentEncodes};
}
