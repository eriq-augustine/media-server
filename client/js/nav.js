"use strict";

var filebrowser = filebrowser || {};
filebrowser.nav = filebrowser.nav || {};

// Start with nothing.
// The hash will be examined before we actually start to override with a location or root.
// Only updateCurrentTarget() is allowed to modify this.
filebrowser.nav._currentTarget = filebrowser.nav._currentTarget || '';
filebrowser.nav._history = filebrowser.nav._history || [];

window.addEventListener("hashchange", function(newValue) {
   if (filebrowser.nav.getCurrentTargetPath() != filebrowser.nav.cleanHashPath()) {
      filebrowser.nav.changeTarget(filebrowser.nav.cleanHashPath());
   }
});

filebrowser.nav.changeTarget = function(path) {
   // TEST
   console.log("Change Target: " + path);

   // Do nothing if we are already pointing to the target.
   // Be careful that we don't block the first load.
   if (filebrowser.nav.getCurrentTargetPath() == path) {
      return;
   }

   var listing = filebrowser.cache.listingFromCache(path);

   if (!listing) {
      filebrowser.cache.loadCache(path, filebrowser.nav.changeTarget.bind(window, path));
      return;
   }

   if (listing.isDir) {
      var files = [];
      $.each(listing.children, function(index, child) {
         files.push(child);
      });
      filebrowser.view.reloadTable(files, path);
   } else {
      filebrowser.view.loadViewer(listing, path);
   }

   // Update the current target.
   filebrowser.nav._updateCurrentTarget(path);
}

filebrowser.nav.getCurrentTargetPath = function() {
   return filebrowser.nav._currentTarget;
}

// This is the only function allowed to modify |_currentTarget|.
filebrowser.nav._updateCurrentTarget = function(path) {
   filebrowser.nav._currentTarget = path;

   // Update the history.
   filebrowser.nav._history.push(path);

   // Change the hash if necessary.
   if (path != filebrowser.nav.cleanHashPath()) {
      window.location.hash = encodeURIComponent(path);
   }
}

// Remove the leading hash and decode the path
filebrowser.nav.cleanHashPath = function() {
   return decodeURIComponent(window.location.hash.replace(/^#/, ''));
}
