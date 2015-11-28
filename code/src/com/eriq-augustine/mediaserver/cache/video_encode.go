package cache;

import (
   "bufio"
   "io/ioutil"
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

var japaneseLangCodes []string = []string{"japanese", "jpn", "jp"};

type AudioStreamSort []map[string]string;

func getEncodePath(cacheDir string) string {
   return filepath.Join(cacheDir, "encode.mp4");
}

func getEncodeDonePath(cacheDir string) string {
   return filepath.Join(cacheDir, "encode.done");
}

func isEncodeDone(cacheDir string) bool {
   return util.PathExists(getEncodeDonePath(cacheDir)) && util.PathExists(getEncodePath(cacheDir));
}

func encodeFileInternal(file model.File, cacheDir string, progressChan chan model.EncodeProgress) error {
   encodePath := getEncodePath(cacheDir);

   // Fetch the info on the streams in this file.
   streamInfo, err := extractStreamInfo(file.DirEntry.Path);
   if (err != nil) {
      return err;
   }

   // Right now we can't deal with multiple audio/video streams so well, so we'll just pick the "best" ones.
   videoStream := getBestVideoStream(streamInfo);
   audioStream := getBestAudioStream(streamInfo);

   encodingThreads := config.GetIntDefault("encodingThreads", 0);

   cmd := exec.Command(
      config.GetString("ffmpegPath"),
      "-i", file.DirEntry.Path, // Input
      "-y", // Overwrite any output files.
      "-nostats",
      "-loglevel", "warning", // Be pretty quiet.
      "-c:v", "libx264", // Video codex.
      "-threads", strconv.Itoa(encodingThreads), // Number of encoding threads to run.
      "-preset", "superfast", // Go as fast as we can without hurting the quality.
      "-vf", "scale=trunc(oh*a/2)*2:720", // Keep the aspect ratio the same, but make the height 720.
                                            // Note that ffmpeg supportd -1 to do this,
                                            // but will sometimes cause odd widths (which mp4 does not allow).
      "-b:v", "1500k", // Experiments show 1500k does well.
      "-maxrate", "1500k", // Should be the bitrate (-b:v).
      "-bufsize", "3000k", // Double the bitrate.
      "-map", "0:" + videoStream["index"],
      "-strict", "-2",
      "-c:a", "aac",
      "-map", "0:" + audioStream["index"],
      "-f", "mp4",
      "-progress", "-", // Send progress to stdout.
      encodePath, // Output.
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
            if (progressChan != nil) {
               progressChan <- model.EncodeProgress{file, cacheDir, currentTimeUS / 1000, videoDurationMS, false};
            }
         }
      }
   }

   err = cmd.Wait()

   if (progressChan != nil) {
      progressChan <- model.EncodeProgress{file, cacheDir, videoDurationMS, videoDurationMS, true};
   }

   if (err != nil) {
      log.ErrorE("Error waiting for encode to finish", err);
   }

   // Mark the encode as complete.
   ioutil.WriteFile(getEncodeDonePath(cacheDir), []byte(""), 0644);

   return nil;
}

// Get the duration in ms.
func getDurrationMS(streamInfo StreamInfo) int64 {
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

   return int64(durationFloat * 1000.0);
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
