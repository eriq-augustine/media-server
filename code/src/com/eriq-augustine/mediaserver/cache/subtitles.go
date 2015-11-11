package cache;

// Handle fetching "posters" from video files.
// A poster is an image to use before the video has played.

import (
   "io/ioutil"
   "os/exec"
   "path/filepath"
   "strconv"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

var subtitleDirs []string = []string{"sub", "subs", "subtitle", "subtitles"}
var subtitleExts []string = []string{"srt", "sub", "sbv", "ass", "ssa", "aqt", "jss", "smi", "vtt", "rt", "pjs", "stl"}

// Get all the subtitles available for |file|.
// This can include internal subtitles (if |file| is a container format),
// adjacent files, and subtitle directories.
func extractSubtitles(file model.File, cacheDir string) ([]string, error) {
   doneFile := filepath.Join(cacheDir, "subtitles.done");

   // Check for the subs before we generate a new one.
   if (util.PathExists(doneFile)) {
      return fetchCachedSubs(cacheDir), nil;
   }

   var numSubtitleFiles = 0;
   numSubtitleFiles = extractSubtitlesFromFile(file.DirEntry.Path, cacheDir, numSubtitleFiles);

   relatedSubFiles := collectRelatedSubFiles(file);

   for _, relatedSubFile := range(relatedSubFiles) {
      numSubtitleFiles = extractSubtitlesFromFile(relatedSubFile, cacheDir, numSubtitleFiles);
   }

   // TODO(eriq): Remove dups

   ioutil.WriteFile(doneFile, []byte(""), 0644);

   return fetchCachedSubs(cacheDir), nil;
}

func fetchCachedSubs(cacheDir string) []string {
   var subs []string = make([]string, 0);

   fileInfos, err := ioutil.ReadDir(cacheDir);
   if (err != nil) {
      log.ErrorE("Unable to read cache dir for subs: " + cacheDir, err);
      return subs;
   }

   for _, fileInfo := range(fileInfos) {
      if (!fileInfo.IsDir() && strings.HasPrefix(fileInfo.Name(), "sub_")) {
         subs = append(subs, filepath.Join(cacheDir, fileInfo.Name()));
      }
   }

   return subs;
}

// Search for subtitle files related to |file|.
func collectRelatedSubFiles(file model.File) []string {
   // Find all the directories to look in (including the current directory).

   // Start with the current directory.
   baseDir := filepath.Dir(file.DirEntry.Path);
   dirs := []string{baseDir};

   // Check for any directories that look like subtitle directories.
   fileInfos, err := ioutil.ReadDir(baseDir);
   if (err == nil) {
      for _, fileInfo := range(fileInfos) {
         if (fileInfo.IsDir() && util.SliceHasString(subtitleDirs, strings.ToLower(fileInfo.Name()))) {
            dirs = append(dirs, filepath.Join(baseDir, fileInfo.Name()));
         }
      }
   }

   // Fetch any subtitles from the directories.
   var subFiles []string = make([]string, 0);
   for _, dir := range(dirs) {
      newSubFiles := fetchSubtitlesFromDirectory(file, dir);
      subFiles = append(subFiles, newSubFiles...);
   }

   return subFiles;
}

func fetchSubtitlesFromDirectory(file model.File, dir string) []string {
   var subs []string = make([]string, 0);
   basename := util.Basename(file.DirEntry.Path);

   // Look for files with the same basename as |file|, but a subtitle extension.
   fileInfos, err := ioutil.ReadDir(dir);
   if (err == nil) {
      for _, fileInfo := range(fileInfos) {
         ext := strings.ToLower(strings.TrimPrefix(fileInfo.Name(), "."));
         suspectBasename := util.Basename(fileInfo.Name());

         if (!fileInfo.IsDir() && basename == suspectBasename && util.SliceHasString(subtitleExts, ext)) {
            subs = append(subs, filepath.Join(dir, fileInfo.Name()));
         }
      }
   }

   return subs;
}

// Get all the available subtitle tracks from a file and put the subs in separate files.
// |nextSubFileIndex| is used to give an easy id to each subtitle file.
// The |nextSubFileIndex| will be maintained in this function and returned back to the caller.
func extractSubtitlesFromFile(path string, cacheDir string, nextSubFileIndex int) int {
   // Fetch the info on the streams in this file.
   streamInfo, err := extractStreamInfo(path);
   if (err != nil) {
      return nextSubFileIndex;
   }

   // Are there any subtitles in this file?
   if (len(streamInfo.Subtitle) == 0) {
      return nextSubFileIndex;
   }

   // ffmpeg can handle multiple subtitle streams at a time, we just need to specify each one.
   args := []string{
      "-i", path,
      "-y", // Overwrite any output files.
      "-nostats",
      "-loglevel", "warning", // Be pretty quiet.
   };

   for _, subStream := range(streamInfo.Subtitle) {
      lang := "und";
      if (util.MapHasKey(subStream, "lang")) {
         lang = subStream["lang"];
      }

      outputFile := filepath.Join(cacheDir, "sub_" + lang + "_" + strconv.Itoa(nextSubFileIndex) + ".vtt");
      streamArgs := []string{
         "-c:s", "webvtt", // Use webvtt for the subtitle format.
         "-map", "0:" + subStream["index"], // Index of the subtitle stream (in the file).
         outputFile,
      }

      args = append(args, streamArgs...);
      nextSubFileIndex++;
   }

   cmd := exec.Command(config.GetString("ffmpegPath"), args...);

   err = cmd.Run();
   if (err != nil) {
      log.ErrorE("Unable to extract subtitles from " + path, err);
      return nextSubFileIndex;
   }

   return nextSubFileIndex;
}
