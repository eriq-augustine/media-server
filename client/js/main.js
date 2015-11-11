"use strict";

var mediaserver = mediaserver || {};

mediaserver.socketPath = 'ws://localhost:1234/ws'
mediaserver.apiPath = 'http://localhost:1234/api/v00/browse/path';
mediaserver.encodeCacheRefreshSec = 10;

// Convert a backend DirEntry to a frontend DirEnt.
function convertBackendDirEntry(dirEntry) {
   if (dirEntry.IsDir) {
      return new filebrowser.Dir(dirEntry.Name, new Date(dirEntry.ModTime));
   } else {
      return new filebrowser.File(dirEntry.Name, new Date(dirEntry.ModTime), dirEntry.Size);
   }
}

// Convert a backend File to a frontend DirEnt.
// Files have more information that just dirents.
function convertBackendFile(file, data) {
   var extraInfo = {
      rawLink: file.RawLink,
      cacheReady: data.CacheReady,
      cacheLink: file.CacheLink,
      poster: file.Poster,
      subtitles: file.Subtitles || []
   };

   return new filebrowser.File(file.DirEntry.Name, new Date(file.DirEntry.ModTime), file.DirEntry.Size, file.RawLink, extraInfo);
}

function fetch(path, callback) {
   path = path || '/';

   var params = {
      "path": path
   };
   var url = mediaserver.apiPath + '?' + $.param(params);

   $.ajax(url, {
      dataType: 'json',
      error: function(request, textStatus, error) {
         // TODO(eriq): log?
         console.log("Error getting data");
         console.log(request);
         console.log(textStatus);
      },
      success: function(data) {
         if (!data.Success) {
            // TODO(eriq): more
            console.log("Unable to get listing");
            console.log(data);
            return;
         }

         var rtnData;
         if (data.IsDir) {
            rtnData = [];
            data.DirEntries.forEach(function(dirEntry) {
               rtnData.push(convertBackendDirEntry(dirEntry));
            });
         } else {
            rtnData = convertBackendFile(data.File, data);
         }

         callback(data.IsDir, rtnData);
      }
   });
}

mediaserver.videoTemplate = `
   <video
      id='main-video-player'
      class='video-player video-js vjs-default-skin vjs-big-play-centered'
   >
      <source src='{{VIDEO_LINK}}' type='{{MIME_TYPE}}'>

      {{SUB_TRACKS}}
      Browser not supported.
   </video>
`;

mediaserver.subtitleTrackTemplate = `
   <track kind="subtitles" src="{{SUB_LINK}}" srclang="{{SUB_LANG}}" label="{{SUB_LABEL}}"></track>
`;

function renderVideo(file) {
   if (!file.extraInfo.cacheReady) {
      return `
         <p>This file needs to be encoded before it can be viewed in-browser.</p>
         <p>After the file is encoded, reload this path.</p>
      `;
   }

   var subTracks = [];
   file.extraInfo.subtitles.forEach(function(sub) {
      var match = sub.match(/sub_(\w+)_(\d+).vtt$/)
      if (!match) {
         return;
      }

      var track = mediaserver.subtitleTrackTemplate;
      track = track.replace('{{SUB_LINK}}', sub);
      track = track.replace('{{SUB_LANG}}', match[1]);
      track = track.replace('{{SUB_LABEL}}', match[1] + '_' + match[2]);

      subTracks.push(track);
   });

   var ext = filebrowser.util.ext(file.extraInfo.cacheLink || file.directLink);
   var mime = '';
   if (filebrowser.filetypes.extensions[ext]) {
      mime = filebrowser.filetypes.extensions[ext].mime;
   }

   var videoHTML = mediaserver.videoTemplate;

   videoHTML = videoHTML.replace('{{VIDEO_LINK}}', file.extraInfo.cacheLink || file.directLink);
   videoHTML = videoHTML.replace('{{MIME_TYPE}}', mime);
   videoHTML = videoHTML.replace('{{SUB_TRACKS}}', subTracks.join());

   return {html: videoHTML, callback: initVideo.bind(this, file)};
}

function initVideo(file) {
   if (videojs.getPlayers()['main-video-player']) {
      videojs.getPlayers()['main-video-player'].dispose();
   }

   videojs('main-video-player', {
      controls: true,
      preload: 'auto',
      poster: file.extraInfo.poster || ''
   });
}

// Look for files that have not encoded yet.
function validateCacheEntry(cacheListing) {
   if (cacheListing.isDir) {
      return true;
   }

   if (cacheListing.extraInfo.cacheReady) {
      return true;
   }

   // If the cache is not ready, give a few seconds between hitting again.
   return ((Date.now() - cacheListing.cacheTime) < (mediaserver.encodeCacheRefreshSec * 1000))
}

$(document).ready(function() {
   // Init the websocket.
   mediaserver.socket.init(mediaserver.socketPath);

   // Init the file browser.
   var options = {
      cacheValidator: validateCacheEntry,
      renderOverrides: {
         video: renderVideo
      }
   };
   filebrowser.init('mediaserver-filebrowser', fetch, options);

   // If there is a valid hash path, follow it.
   // Otherwise, set up a new hash at root.
   var target = '/';
   if (window.location.hash) {
      target = filebrowser.nav.cleanHashPath();
   }

   filebrowser.nav.changeTarget(target);
});
