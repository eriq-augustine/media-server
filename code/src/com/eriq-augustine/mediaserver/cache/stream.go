package cache;

// Extract the stream information from a video file using ffprobe.

import (
   "os/exec"
   "regexp"
   "strings"

   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/util"
)

const (
   PARSE_PROBE_STATE_OPEN = iota
   PARSE_PROBE_STATE_METADATA
   PARSE_PROBE_STATE_STREAM
)

type StreamInfo struct {
   Metadata map[string]string
   Video []map[string]string
   Audio []map[string]string
   Subtitle []map[string]string
   Other []map[string]string
}

func NewStreamInfo() StreamInfo {
   var info StreamInfo;
   info.Metadata = make(map[string]string);
   info.Video = make([]map[string]string, 0);
   info.Audio = make([]map[string]string, 0);
   info.Subtitle = make([]map[string]string, 0);
   info.Other = make([]map[string]string, 0);
   return info;
}

func extractStreamInfo(path string) (StreamInfo, error) {
   cmd := exec.Command(
      config.GetString("ffprobePath"),
      "-hide_banner",
      "-show_streams",
      "-show_format",
      path,
   );

   output, err := cmd.Output();
   if (err != nil) {
      log.ErrorE("Unable to extract streams", err);
      return NewStreamInfo(), err;
   }

   streamInfo := parseProbe(string(output));

   return streamInfo, nil;
}

func parseProbe(output string) StreamInfo {
   var streamInfo StreamInfo = NewStreamInfo();
   var currentStream map[string]string = make(map[string]string);

   state := PARSE_PROBE_STATE_OPEN
   lines := regexp.MustCompile(`\r?\n`).Split(output, -1);

   for _, line := range(lines) {
      line = strings.TrimSpace(line);

      switch state {
      case PARSE_PROBE_STATE_OPEN:
         if (line == "[FORMAT]") {
            state = PARSE_PROBE_STATE_METADATA;
         } else if (line == "[STREAM]") {
            state = PARSE_PROBE_STATE_STREAM;
            currentStream = make(map[string]string);
         }
         break;
      case PARSE_PROBE_STATE_METADATA:
         if (line == "[/FORMAT]") {
            state = PARSE_PROBE_STATE_OPEN;
         } else {
            data := strings.SplitN(strings.TrimPrefix(strings.ToLower(line), "tag:"), "=", 2)
            if (strings.TrimSpace(data[1]) != "") {
               streamInfo.Metadata[strings.TrimSpace(data[0])] = strings.TrimSpace(data[1]);
            }
         }
         break;
      case PARSE_PROBE_STATE_STREAM:
         if (line == "[/STREAM]") {
            switch currentStream["codec_type"] {
            case "video":
               streamInfo.Video = append(streamInfo.Video, currentStream);
            case "audio":
               streamInfo.Audio = append(streamInfo.Audio, currentStream);
            case "subtitle":
               // Ensure that "lang" is populated if "language" exists.
               if (util.MapHasKey(currentStream, "language") && !util.MapHasKey(currentStream, "lang")) {
                  currentStream["lang"] = currentStream["language"];
               }

               streamInfo.Subtitle = append(streamInfo.Subtitle, currentStream);
            case "attachment":
               streamInfo.Other = append(streamInfo.Other, currentStream);
            default:
               log.Warn("Unknown codec_type: " + currentStream["codec_type"] + ".");
               streamInfo.Other = append(streamInfo.Other, currentStream);
            }

            currentStream = make(map[string]string);
            state = PARSE_PROBE_STATE_OPEN;
         } else {
            data := strings.SplitN(strings.TrimPrefix(strings.ToLower(line), "tag:"), "=", 2)
            if (strings.TrimSpace(data[1]) != "") {
               currentStream[strings.TrimSpace(data[0])] = strings.TrimSpace(data[1]);
            }
         }
         break;
      default:
         log.Fatal("Unknown parse state: " + string(state));
         break;
      }
   }

   return streamInfo;
}
