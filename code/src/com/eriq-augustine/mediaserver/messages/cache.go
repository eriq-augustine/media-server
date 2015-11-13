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
   // Clear out any potentially sentitive information.
   if (progress != nil) {
      // Make a copy so we don't mess with original.
      var progressCopy model.EncodeProgress = *progress;
      progressCopy.CacheDir = "";
      progress = &progressCopy;
   }

   for _, request := range(queue) {
      request.CacheDir = "";
   }

   return &CacheStatus{true, progress, queue, recentEncodes};
}
