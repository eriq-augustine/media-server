package cache;

import (
   "bufio"
   "fmt"
   "os/exec"
   "path/filepath"
   "sort"
   "strconv"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

const (
   REQUEST_BUFFER_SIZE = 1024
)

var japaneseLangCodes []string = []string{"japanese", "jpn", "jp"};
var requestQueue chan EncodeRequest;

type EncodeRequest struct {
   File model.File
   CacheDir string
}

type AudioStreamSort []map[string]string;

// Recieve encode requests and perform them one at a time.
func manageEncodesThread() {
   // Loop for every on the channel.
   for request := range(requestQueue) {
      // TEST
      log.Debug("Request: " + request.CacheDir);

      encodeFile(request.File, request.CacheDir);

      // TEST
      log.Debug("Done: " + request.CacheDir);
   }
}

func getEncodePath(cacheDir string) string {
   return filepath.Join(cacheDir, "encode.mp4");
}

// The second returned value indicates if the encode is good.
func requestEncode(file model.File, cacheDir string) (string, bool) {
   // First see if this is the first request and we need to start the encoding thread.
   if (requestQueue == nil) {
      requestQueue = make(chan EncodeRequest, REQUEST_BUFFER_SIZE);
      go manageEncodesThread();
   }

   encodePath := getEncodePath(cacheDir);

   // Check for the encode before we generate a new one.
   if (util.PathExists(encodePath)) {
      return encodePath, true;
   }

   // Pass on the request.
   requestQueue <- EncodeRequest{file, cacheDir};

   return "", false;
}

func encodeFile(file model.File, cacheDir string) error {
   encodePath := getEncodePath(cacheDir);

   // Fetch the info on the streams in this file.
   streamInfo, err := extractStreamInfo(file.DirEntry.Path);
   if (err != nil) {
      return err;
   }

   // Right now we can't deal with multiple audio/video streams so well, so we'll just pick the "best" ones.
   videoStream := getBestVideoStream(streamInfo);
   audioStream := getBestAudioStream(streamInfo);

   encodingThreads := config.GetIntDefault("encodingThreads", 2);

   cmd := exec.Command(
      config.GetString("ffmpegPath"),
      "-i", file.DirEntry.Path, // Input
      "-y", // Overwrite any output files.
      "-nostats",
      "-loglevel", "warning", // Be pretty quiet.
      "-c:v", "libx264", // Video codex.
      "-crf", "28", // Constant rate factor.
      "-threads", strconv.Itoa(encodingThreads), // Number of encoding threads to run.
      "-preset", "veryfast", // Go as fast as we can.
      "-b:v", util.MapGetWithDefault(streamInfo.Metadata, "bit_rate", "2426222"),
      "-map", "0:" + videoStream["index"],
      "-strict", "-2",
      "-c:a", "aac",
      "-map", "0:" + audioStream["index"],
      "-progress", "-", // Send progress to stdout.
      encodePath, // Output. File format is infered from extension.
   );

   stdout, err := cmd.StdoutPipe()
   if err != nil {
      log.ErrorE("Unable to get stdout from encode process", err);
      return err;
   }

   err = cmd.Start();
   if (err != nil) {
      log.ErrorE("Unable to encode file", err);
      return err;
   }

   videoDurationMS := getDurrationMS(streamInfo);

   // Read all the progress information.
   bufferedStdout := bufio.NewReader(stdout);
   for {
      line, _, err := bufferedStdout.ReadLine();
      if (err != nil || line == nil) {
         break;
      }

      // Even though is has the suffix "ms", this is really in microseconds, not milliseconds.
      if (strings.HasPrefix(string(line), "out_time_ms")) {
         currentTimeUS, err := strconv.ParseInt(strings.TrimPrefix(string(line), "out_time_ms="), 10, 64);
         if (err == nil) {
            // TEST
            if (1 == 2) {
               fmt.Printf("%d / %d\n", currentTimeUS / 1000, videoDurationMS);
            }
         }
      }
   }

   err = cmd.Wait()
   if (err != nil) {
      log.ErrorE("Error waiting for encode to finish", err);
   }

   return nil;
}

// Get the duration in ms.
func getDurrationMS(streamInfo StreamInfo) int {
   durationString := "-1";

   // First check the metadata, then the video and audio streams.
   streamsToCheck := []map[string]string{streamInfo.Metadata};
   streamsToCheck = append(streamsToCheck, streamInfo.Video...);
   streamsToCheck = append(streamsToCheck, streamInfo.Audio...);
   streamsToCheck = append(streamsToCheck, streamInfo.Subtitle...);
   streamsToCheck = append(streamsToCheck, streamInfo.Other...);

   for _, stream := range(streamsToCheck) {
      if (util.MapHasKey(stream, "duration")) {
         durationString = stream["duration"];
         break;
      }
   }

   durationFloat, err := strconv.ParseFloat(durationString, 64);
   if (err != nil) {
      return -1;
   }

   return int(durationFloat * 1000.0);
}

// For video streams, just pick the first one.
func getBestVideoStream(streamInfo StreamInfo) map[string]string {
   if (len(streamInfo.Video) == 0) {
      return nil;
   }

   return streamInfo.Video[0];
}

// For audio streams, first remove ant commentary and then favor japanese.
func getBestAudioStream(streamInfo StreamInfo) map[string]string {
   if (len(streamInfo.Audio) == 0) {
      return nil;
   }

   sort.Sort(AudioStreamSort(streamInfo.Audio));

   return streamInfo.Audio[0];
}

func (arr AudioStreamSort) Len() int {
   return len(arr);
}

func (arr AudioStreamSort) Swap(i int, j int) {
    arr[i], arr[j] = arr[j], arr[i];
}

func (arr AudioStreamSort) Less(i int, j int) bool {
   // Check commentary first.
   iHasCommentary := strings.Contains(strings.ToLower(util.MapGetWithDefault(arr[i], "title", "")), "commentary");
   jHasCommentary := strings.Contains(strings.ToLower(util.MapGetWithDefault(arr[j], "title", "")), "commentary");

   if (iHasCommentary != jHasCommentary) {
      return jHasCommentary;
   }

   // Then check language (favor japanese).
   iHasJpn := util.SliceHasString(japaneseLangCodes, strings.ToLower(util.MapGetWithDefault(arr[i], "lang", "")));
   jHasJpn := util.SliceHasString(japaneseLangCodes, strings.ToLower(util.MapGetWithDefault(arr[j], "lang", "")));

   if (iHasJpn != jHasJpn) {
      return iHasJpn;
   }

   // Finally, just use the stream index.
   iIndex := strings.ToLower(util.MapGetWithDefault(arr[i], "index", "99"));
   jIndex := strings.ToLower(util.MapGetWithDefault(arr[j], "index", "99"));

   return iIndex < jIndex;
}
