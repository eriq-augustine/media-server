"use strict";

var filebrowser = filebrowser || {};
filebrowser.nav = filebrowser.nav || {};

// Start with nothing.
// The hash will be examined before we actually start to override with a location or root.
// Only updateCurrentTarget() is allowed to modify this.
filebrowser.nav._currentTarget = filebrowser.nav._currentTarget || '';
filebrowser.nav._history = filebrowser.nav._history || [];

filebrowser.DirEnt = function(name, modDate, size, isDir, detailsCached) {
   this.name = name;
   this.modDate = modDate;
   this.size = size;
   this.isDir = isDir;
   this.detailsCached = detailsCached;
}

filebrowser.Dir = function(name, modDate) {
   // TODO(eriq): Figure out the best time when none is supplied. Now? Zero time?
   modDate = modDate || new Date();

   filebrowser.DirEnt.call(this, name, modDate, 0, true, false);
   this.children = {};
}

filebrowser.Dir.prototype = Object.create(filebrowser.DirEnt.prototype);
filebrowser.Dir.prototype.constructor = filebrowser.Dir;

filebrowser.File = function(name, modDate, size, directLink, extraInfo) {
   extraInfo = extraInfo || {};

   filebrowser.DirEnt.call(this, name, modDate, size, false, false);
   this.directLink = directLink;
   this.extraInfo = extraInfo;

   if (name.includes('.')) {
      var nameParts = name.match(/^(.*)\.([^\.]*)$/);
      this.basename = nameParts[1];
      this.extension = nameParts[2];
   } else {
      this.basename = name;
      this.extension = '';
   }
}

filebrowser.File.prototype = Object.create(filebrowser.DirEnt.prototype);
filebrowser.File.prototype.constructor = filebrowser.File;

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

function arrayToTableRow(data, isHeader) {
   isHeader = isHeader | false;
   var cellType = isHeader ? 'th' : 'td';

   var tr = document.createElement('tr');

   data.forEach(function(dataObject) {
      var td = document.createElement(cellType);
      td.appendChild(document.createTextNode(dataObject));
      tr.appendChild(td);
   });

   return tr;
}

function fileToTableRow(file) {
   // TODO(eriq): MIME types
   var data = [file.name, file.modDate, 'TODO', file.size];
   return arrayToTableRow(data, false);
}

function joinURL(base, addition) {
   if (!base || base == '/') {
      return '/' + addition;
   }

   return base + '/' + addition;
}

function filesToTable(path, files) {
   var table = document.createElement('table');

   var tableHead = document.createElement('thead');
   var headerData = ['Name', 'Date', 'Type', 'Size'];
   tableHead.appendChild(arrayToTableRow(headerData, true));
   table.appendChild(tableHead);

   var tableBody = document.createElement('tbody');
   files.forEach(function(file) {
      var row = fileToTableRow(file);
      var url = joinURL(path, file.name);
      row.setAttribute('data-path', url);
      row.addEventListener('click', listingClicked.bind(window, file, url));
      tableBody.appendChild(row);
   });
   table.appendChild(tableBody);

   return table;
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

function listingClicked(file, path, ev) {
   // TEST
   console.log("Listing clicked: " + path);
   console.log(file);

   changeTarget(path);
}

function changeTarget(path) {
   // TEST
   console.log("Change Target: " + path);

   // Do nothing if we are already pointing to the target.
   // Be careful that we don't block the first load.
   if (getCurrentTargetPath() == path) {
      return;
   }

   var listing = filebrowser.cache.listingFromCache(path);

   if (!listing) {
      filebrowser.cache.loadCache(path, changeTarget.bind(window, path));
      return;
   }

   if (listing.isDir) {
      var files = [];
      $.each(listing.children, function(index, child) {
         files.push(child);
      });
      reloadTable(files, path);
   } else {
      loadViewer(listing, path);
   }

   // Update the current target.
   updateCurrentTarget(path);
}

function getCurrentTargetPath() {
   return filebrowser.nav._currentTarget;
}

// This is the only function allowed to modify |_currentTarget|.
function updateCurrentTarget(path) {
   filebrowser.nav._currentTarget = path;

   // Update the history.
   filebrowser.nav._history.push(path);

   // Change the hash if necessary.
   if (path != cleanHashPath()) {
      window.location.hash = encodeURIComponent(path);
   }
}

function loadViewer(file, path) {
   console.log('loadViewer');
   console.log(file);

   // TODO(eriq): Re-architect the html some, it's not just a table.
   clearTable();

   $('#tableArea').html(filebrowser.filetypes.renderHTML(file));
}

function reloadTable(files, path) {
   var table = filesToTable(path, files);
   table.id = 'myTable';
   table.className = 'tablesorter';

   // TODO(eriq): Better ids
   clearTable();
   $('#tableArea').append(table);

   $("#myTable").tablesorter({
      sortList: [[0,0]],
      widgets: ['zebra']
   });
}

function clearTable() {
   $('#tableArea').empty();
}

// Remove the leading hash and decode the path
function cleanHashPath() {
   return decodeURIComponent(window.location.hash.replace(/^#/, ''));
}

window.addEventListener("hashchange", function(newValue) {
   if (getCurrentTargetPath() != cleanHashPath()) {
      changeTarget(cleanHashPath());
   }
});

$(document).ready(function() {
   // If there is a valid hash path, follow it.
   // Otherwise, set up a new hash at root.
   var target = '/';
   if (window.location.hash) {
      target = cleanHashPath();
   }

   changeTarget(target);
});
