package cache;

// Code for requesting encodes, keeping track of the progress, and recently encoded files.

import (
   "sort"
   "sync"
   "time"

   "com/eriq-augustine/mediaserver/model"
)

const (
   DEFAULT_RECENT_ENCODE_COUNT = 20
   REQUEST_BUFFER_SIZE = 1024
)

// This will not block.
// Once this is called, the manager owns the encode.
func queueEncode(file model.File, cacheDir string) {
   encodeMutex.Lock();
   defer encodeMutex.Unlock();

   _, exists := allEncodeRequests[cacheDir];
   if (exists) {
      return;
   }
   allEncodeRequests[cacheDir] = true;

   queue = append(queue, model.EncodeRequest{file, cacheDir, time.Now()});
   startNextEncode();
}

func getEncodeStatus() model.EncodeStatus {
   encodeMutex.Lock();
   defer encodeMutex.Unlock();

   var progressCopy *model.EncodeProgress = nil;
   if (progress != nil) {
      temp := *progress;
      progressCopy = &temp;
   }

   encodeCount := DEFAULT_RECENT_ENCODE_COUNT;
   if (encodeCount > len(complete)) {
      encodeCount = len(complete);
   }

   return model.EncodeStatus{
      Queue: queue,
      Progress: progressCopy,
      Complete: complete[0:encodeCount],
   };
}

// Add complete encodes to the manager.
// These will typically come from scanning the cache after initialization.
func loadRecentVideoEncodes(encodes []model.CompleteEncode) {
   encodeMutex.Lock();
   defer encodeMutex.Unlock();

   for _, encode := range(encodes) {
      allEncodeRequests[encode.CacheDir] = true;
   }

   complete = encodes;
   sort.Sort(ByCompleteTime(complete));
}

func removeEncode(cacheDir string) {
   encodeMutex.Lock();
   defer encodeMutex.Unlock();

   delete(allEncodeRequests, cacheDir);

   // Check to see if it is in the list of complete encodes.
   for i, encode := range(complete) {
      if (encode.CacheDir == cacheDir) {
         complete = append(complete[:i], complete[i + 1:]...);
         break;
      }
   }
}

// Private

var encodeMutex *sync.Mutex;

// Keep track of all encode requests so we don't double queue.
// {cacheDir: true}
var allEncodeRequests map[string]bool;

// All of these should be protected with |encodeMutex|.
var queue []model.EncodeRequest
var inProgress *model.EncodeRequest
var progress *model.EncodeProgress
var complete []model.CompleteEncode

var nextEncodeChan chan model.EncodeRequest;
var progressChan chan model.EncodeProgress;

func init() {
   encodeMutex = &sync.Mutex{};

   allEncodeRequests = make(map[string]bool);

   queue = make([]model.EncodeRequest, 0);
   inProgress = nil;
   progress = nil;
   complete = make([]model.CompleteEncode, 0);

   nextEncodeChan = make(chan model.EncodeRequest, REQUEST_BUFFER_SIZE);
   progressChan = make(chan model.EncodeProgress, REQUEST_BUFFER_SIZE);

   go encoderThread();
   go progressThread();
}

func encoderThread() {
   // Wait for requests forever.
   for request := range(nextEncodeChan) {
      encodeFileInternal(request.File, request.CacheDir, progressChan);
      encodeComplete();
   }
}

// Watch for encode progress.
func progressThread() {
   // Wait for requests forever.
   for update := range(progressChan) {
      updateProgress(update);
   }
}

func updateProgress(update model.EncodeProgress) {
   encodeMutex.Lock();
   defer encodeMutex.Unlock();

   // It is possible for the progress thread to fall behind,
   // so check the conditions before making updates.
   if (inProgress == nil) {
      return;
   }

   // Make sure that the update is for what we think is in progress.
   if (inProgress.CacheDir != update.CacheDir) {
      return;
   }

   progress = &update;
}

// Note that the caller should have |encodeMutex| locked before calling.
func startNextEncode() {
   if (len(queue) == 0 || inProgress != nil) {
      return;
   }

   nextEncode := queue[0];

   queue = queue[1:];
   inProgress = &nextEncode;

   nextEncodeChan <- nextEncode;
}

func encodeComplete() {
   encodeMutex.Lock();
   defer encodeMutex.Unlock();

   cacheDir := inProgress.CacheDir;
   completeEncode := model.CompleteEncode{inProgress.File, time.Now(), getEncodePath(cacheDir), cacheDir};

   // Settle the finished encode.
   // Make sure to prepend the encode so it does not get cut out when we
   // choose only the most recent encodes to send to the users.
   complete = append([]model.CompleteEncode{completeEncode}, complete...);
   inProgress = nil;
   progress = nil;

   setCachedVideoEncode(cacheDir, completeEncode);

   startNextEncode();
}

// Sort complete encodes by time.
type ByCompleteTime []model.CompleteEncode;

func (this ByCompleteTime) Len() int { return len(this); }
func (this ByCompleteTime) Swap(i int, j int) { this[i], this[j] = this[j], this[i]; }
func (this ByCompleteTime) Less(i int, j int) bool { return this[i].CompleteTime.Before(this[j].CompleteTime); }
