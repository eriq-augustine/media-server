"use strict";

var mediaserver = mediaserver || {};

mediaserver.apiPath = '/api/v00/browse/path';
mediaserver.socketPath = 'ws://' + window.location.host + '/ws';
mediaserver.encodeCacheRefreshSec = 10;

// Convert a backend DirEntry to a frontend DirEnt.
mediaserver._convertBackendDirEntry = function(dirEntry) {
   if (dirEntry.IsDir) {
      return new filebrowser.Dir(dirEntry.Name, new Date(dirEntry.ModTime));
   } else {
      return new filebrowser.File(dirEntry.Name, new Date(dirEntry.ModTime), dirEntry.Size);
   }
}

// Convert a backend File to a frontend DirEnt.
// Files have more information that just dirents.
mediaserver._convertBackendFile = function(file, data) {
   var extraInfo = {
      rawLink: file.RawLink,
      cacheReady: data.CacheReady,
      cacheLink: file.CacheLink,
      poster: file.Poster,
      subtitles: file.Subtitles || []
   };

   return new filebrowser.File(file.DirEntry.Name, new Date(file.DirEntry.ModTime), file.DirEntry.Size, file.RawLink, extraInfo);
}

mediaserver._fetch = function(path, callback) {
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
               rtnData.push(mediaserver._convertBackendDirEntry(dirEntry));
            });
         } else {
            rtnData = mediaserver._convertBackendFile(data.File, data);
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

mediaserver.renderEncodeActivity = function(encodeActivity) {
   // If there is no activity, just the encode activity area.
   if (!encodeActivity.Progress && encodeActivity.Queue.length == 0 && encodeActivity.RecentEncodes.length == 0) {
      $('.encode-activity-container .current-encode').empty();
      $('.encode-activity-container .queue').empty();
      $('.encode-activity-container .recent-encodes').empty();
      $('.encode-activity-container').hide();
      return;
   }

   var encoding = document.createElement('div');
   encoding.className = 'encoding';
   if (encodeActivity.Progress) {
      var encodeActivityElement = mediaserver._renderEncodeActivityItem(encodeActivity.Progress.File);

      var progressBar = document.createElement('progress');
      progressBar.setAttribute('max', encodeActivity.Progress.TotalMS);
      progressBar.setAttribute('value', encodeActivity.Progress.CompleteMS);
      encodeActivityElement.appendChild(progressBar);

      encoding.appendChild(encodeActivityElement);
   } else {
     var nothingEncoding = document.createElement('span');
      nothingEncoding.className = 'encode-list-element-placebolder';
      nothingEncoding.textContent = 'Nothing Encoding';
      encoding.appendChild(nothingEncoding);
   }

   var queue = document.createElement('div');
   queue.className = 'queue';
   if (encodeActivity.Queue.length > 0) {
      encodeActivity.Queue.forEach(function(queueItem) {
         queue.appendChild(mediaserver._renderEncodeActivityItem(queueItem.File));
      });
   } else {
      var nothingQueued = document.createElement('span');
      nothingQueued.className = 'encode-list-element-placebolder';
      nothingQueued.textContent = 'Nothing Queued';
      queue.appendChild(nothingQueued);
   }

   var recentlyEncoded = document.createElement('div');
   recentlyEncoded.className = 'recently-encoded';
   if (encodeActivity.RecentEncodes.length > 0) {
      encodeActivity.RecentEncodes.forEach(function(recentEncode) {
         recentlyEncoded.appendChild(mediaserver._renderEncodeActivityItem(recentEncode.File));
      });
   } else {
      var nothingRecent = document.createElement('span');
      nothingRecent.className = 'encode-list-element-placebolder';
      nothingRecent.textContent = 'No Recent Encodes';
      recentlyEncoded.appendChild(nothingRecent);
   }

   $('.encode-activity-container .current-encode').empty().append(encoding);
   $('.encode-activity-container .queue').empty().append(queue);
   $('.encode-activity-container .recent-encodes').empty().append(recentlyEncoded);
   $('.encode-activity-container').show();
}

mediaserver._renderEncodeActivityItem = function(file) {
   var encodeActivityItem = document.createElement('div');
   encodeActivityItem.className = 'encode-list-element';
   encodeActivityItem.addEventListener('click', filebrowser.nav.changeTarget.bind(window, '/' + file.DirEntry.AbstractPath));

   var fileName = document.createElement('span');
   fileName.textContent = file.DirEntry.Name;

   encodeActivityItem.appendChild(fileName);
   return encodeActivityItem;
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

mediaserver._renderVideo = function(file) {
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

   return {html: videoHTML, callback: mediaserver._initVideo.bind(this, file)};
}

mediaserver._initVideo = function(file) {
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
mediaserver._validateCacheEntry = function(cacheListing) {
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
      cacheValidator: mediaserver._validateCacheEntry,
      renderOverrides: {
         video: mediaserver._renderVideo
      }
   };
   filebrowser.init('mediaserver-filebrowser', mediaserver._fetch, options);

   // If there is a valid hash path, follow it.
   // Otherwise, set up a new hash at root.
   var target = '/';
   if (window.location.hash) {
      target = filebrowser.nav.cleanHashPath();
   }

   filebrowser.nav.changeTarget(target);
});
