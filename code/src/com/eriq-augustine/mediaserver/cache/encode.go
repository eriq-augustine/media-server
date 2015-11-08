package cache;

import (
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

func encodeFile(file model.File, cacheDir string) error {
   encodePath := filepath.Join(cacheDir, "encode.mp4");

   // Check for the encode before we generate a new one.
   if (util.PathExists(encodePath)) {
      return nil;
   }

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
      encodePath, // Output. File format is infered from extension.
   );

   err = cmd.Run();
   if (err != nil) {
      log.ErrorE("Unable to encode file", err);
      return err;
   }

   return nil;
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
