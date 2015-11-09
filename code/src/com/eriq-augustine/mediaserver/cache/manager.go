package cache;

// Code for requesting encodes, keeping track of the progress, and recently encoded files.

import (
   "fmt"

   "com/eriq-augustine/mediaserver/model"
)

const (
   DEFAULT_RECENT_ENCODE_COUNT = 20
)

// TODO(eriq): Consider doing updates with pointers insead since we will be sending many.

var encodeRequestChan chan EncodeRequest;
var manager EncodeManager;

type EncodeRequest struct {
   File model.File
   CacheDir string
}

type EncodeProgress struct {
   File model.File
   CacheDir string
   CompleteMS int64
   TotalMS int64
   Done bool
}

type EncodeManager struct {
   Queue []EncodeRequest
   InProgress *EncodeRequest
   Progress *EncodeProgress
   Complete []EncodeRequest
   ProgressChan chan EncodeProgress
   NextEncodeChan chan EncodeRequest
}

// Setup the proper threads.
func init() {
   // Set up the chans.
   // This is an external channel to make requests to the manager.
   encodeRequestChan = make(chan EncodeRequest, REQUEST_BUFFER_SIZE);

   // These chans are only between the manager and encoder threads.
   manager.ProgressChan  = make(chan EncodeProgress, REQUEST_BUFFER_SIZE);
   manager.NextEncodeChan = make(chan EncodeRequest, REQUEST_BUFFER_SIZE);

   // Set up the manager.
   manager.Queue = make([]EncodeRequest, 0);
   manager.InProgress = nil;
   manager.Progress = nil;
   manager.Complete = make([]EncodeRequest, 0);

   // Set up the threads.
   go encoderThread(manager.NextEncodeChan, manager.ProgressChan);
   go managerThread(&manager);
}

func encoderThread(nextEncodeChan chan EncodeRequest, progressChan chan EncodeProgress) {
   // Wait for requests forever.
   for request := range(nextEncodeChan) {
      encodeFileInternal(request.File, request.CacheDir, progressChan);
   }
}

// Recieve encode requests and perform them one at a time.
func managerThread(manager *EncodeManager) {
   var request EncodeRequest;
   var progress EncodeProgress;

   // Select on encode requests and progress updates.
   for {
      select {
      case request = <- encodeRequestChan:
         manager.QueueRequest(request);
         break;
      case progress = <- manager.ProgressChan:
         manager.UpdateProgress(progress);
         break;
      }
   }
}

func (manager *EncodeManager) UpdateProgress(update EncodeProgress) {
   if (update.Done) {
      manager.encodeComplete();
   } else {
      manager.Progress = &update;
   }
}

func (manager *EncodeManager) QueueRequest(request EncodeRequest) {
   manager.Queue = append(manager.Queue, request);
   manager.startNextEncode();
}

func (manager *EncodeManager) encodeComplete() {
   // Settle the finished encode.
   manager.Complete = append(manager.Complete, *manager.InProgress);
   manager.InProgress = nil;
   manager.Progress = nil;

   // Setup the next one.
   manager.startNextEncode();
}

func (manager *EncodeManager) startNextEncode() {
   if (len(manager.Queue) == 0 || manager.InProgress != nil) {
      return;
   }

   nextEncode := manager.Queue[0];

   manager.Queue = manager.Queue[1:];
   manager.InProgress = &nextEncode;

   manager.NextEncodeChan <- nextEncode;
}

func getProgress() *EncodeProgress {
   return manager.Progress;
}

func getQueue() []EncodeRequest {
   return manager.Queue;
}

func getRecentEncodes(count int) []EncodeRequest {
   if (count <= 0) {
      count = DEFAULT_RECENT_ENCODE_COUNT;
   }

   if (count > len(manager.Complete)) {
      return manager.Complete;
   }

   // TODO(eriq): This is a little unsafe.
   return manager.Complete[0:count];
}

// This will not block.
// Once this is called, the manager owns the encode.
func queueEncode(file model.File, cacheDir string) {
   // TEST
   fmt.Println(cacheDir);

   // Pass on the request.
   encodeRequestChan <- EncodeRequest{file, cacheDir};
}
