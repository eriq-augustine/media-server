"use strict";

// Convert a backend DirEntry to a frontend DirEnt.
function convertBackendDirEntry(dirEntry) {
   if (dirEntry.IsDir) {
      return new filebrowser.Dir(dirEntry.Name, dirEntry.ModTime);
   } else {
      return new filebrowser.File(dirEntry.Name, dirEntry.ModTime, dirEntry.Size);
   }
}

// Convert a backend File to a frontend DirEnt.
// Files have more information that just dirents.
function convertBackendFile(file) {
   var extraInfo = {
      cacheLink: file.CacheLink,
      rawLink: file.RawLink
   };

   return new filebrowser.File(file.DirEntry.Name, file.DirEntry.ModTime, file.DirEntry.Size, file.RawLink, extraInfo);
}

function fetch(path, callback) {
   path = path || '/';

   var params = {
      "path": path
   };
   var url = 'http://localhost:1234/api/v00/browse/path?' + $.param(params);

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

         // TEST
         console.log(data);

         var rtnData;
         if (data.IsDir) {
            rtnData = [];
            data.DirEntries.forEach(function(dirEntry) {
               rtnData.push(convertBackendDirEntry(dirEntry));
            });
         } else {
            rtnData = convertBackendFile(data.File);
         }

         callback(data.IsDir, rtnData);
      }
   });
}

$(document).ready(function() {
   // Register the function for fetching files from the server.
   filebrowser.init('mediaserver-filebrowser', fetch);

   // If there is a valid hash path, follow it.
   // Otherwise, set up a new hash at root.
   var target = '/';
   if (window.location.hash) {
      target = filebrowser.nav.cleanHashPath();
   }

   filebrowser.nav.changeTarget(target);
});
