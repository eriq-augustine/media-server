package cache;

// Handle fetching "posters" from video files.
// A poster is an image to use before the video has played.

import (
   "os/exec"
   "path/filepath"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/model"
   "com/eriq-augustine/mediaserver/util"
)

const (
   POSTER_TIME_SEC = "20"
)

func fetchPoster(file *model.File, cacheDir string) (*string, error) {
   posterPath := filepath.Join(cacheDir, "poster.png");

   // Check for the poster before we generate a new one.
   if (util.PathExists(posterPath)) {
      return &posterPath, nil;
   }

   cmd := exec.Command(
      config.GetString("ffmpegPath"),
      "-i", file.DirEntry.Path, // Input
      "-y", // Overwrite any output files.
      "-an", // Don't do any audio.
      "-nostats",
      "-loglevel", "warning", // Be pretty quiet.
      "-ss", POSTER_TIME_SEC, // The time to take the screenshot (in seconds).
      "-vframes", "1", // Take only one frame (screenshot).
      posterPath, // Output. File format is infered from extension.
   );

   err := cmd.Run();
   if (err != nil) {
      log.ErrorE("Unable to generate poster", err);
      return nil, err;
   }

   return &posterPath, nil;
}
