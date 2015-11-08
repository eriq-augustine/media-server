package cache;

// Handle fetching "posters" from video files.
// A poster is an image to use before the video has played.

import (
   "com/eriq-augustine/mediaserver/model"
)

func extractSubtitles(file model.File, cacheDir string) {
   // TEST
   extractStreamInfo(file);
}

/*TEST
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

SUBTITLE_DIRS = ['sub', 'subs', 'subtitle', 'subtitles']
SUBTITLE_EXTS = ['srt', 'sub', 'sbv', 'ass', 'ssa', 'aqt', 'jss', 'smi', 'vtt', 'rt', 'pjs', 'stl']

   doneFile := filepath.Join(cacheDir, "subtitles.done");

   // Check for the subs before we generate a new one.
   if (util.PathExists(doneFile)) {
      return nil;
   }

func fetchPoster(file model.File, cacheDir string) error {
   posterPath := filepath.Join(cacheDir, "poster.png");

   // Check for the poster before we generate a new one.
   if (util.PathExists(posterPath)) {
      return nil;
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
      return err;
   }

   return nil;
}
*/
