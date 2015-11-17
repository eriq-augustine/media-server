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

filebrowser.nav.changeTarget = function(path, count) {
   // Do nothing if we are already pointing to the target.
   // Be careful that we don't block the first load.
   if (filebrowser.nav.getCurrentTargetPath() == path) {
      return;
   }

   var listing = filebrowser.cache.listingFromCache(path);

   if (!listing) {
      filebrowser.cache.loadCache(path, filebrowser.nav.changeTarget.bind(window, path, count + 1));
      return;
   }

   if (listing.isDir) {
      var files = [];
      $.each(listing.children, function(index, child) {
         files.push(child);
      });
      filebrowser.view.loadBrowserContent(listing, files, path);
   } else {
      filebrowser.view.loadViewer(listing, path);
   }

   // Update the current target.
   filebrowser.nav._updateCurrentTarget(path, listing);
}

filebrowser.nav.getCurrentTargetPath = function() {
   return filebrowser.nav._currentTarget;
}

// This is the only function allowed to modify |_currentTarget|.
filebrowser.nav._updateCurrentTarget = function(path, listing) {
   filebrowser.nav._currentTarget = path;

   // Update the history.
   filebrowser.nav._history.push(path);

   // Change the hash if necessary.
   if (path != filebrowser.nav.cleanHashPath()) {
      window.location.hash = filebrowser.nav.encodeForHash(path);
   }

   // TODO(eriq): Whether or not to change the title (and hash) should be an option.
   // Change the page's title.
   document.title = filebrowser.util.basename(path);

   // Update the breadcrumbs.
   filebrowser.view.loadBreadcrumbs(filebrowser.nav._buildBreadcrumbs(path));

   // Update any context actions.
   filebrowser.view.loadContextActions(listing, path);
}

filebrowser.nav._buildBreadcrumbs = function(path) {
   var breadcrumbs = [];

   var runningPath = '';
   var pathArray = path.replace(/\/$/, '').split('/');
   pathArray.forEach(function(pathElement) {
      // Replace the first element (empty) with root.
      pathElement = pathElement || '/';

      runningPath = filebrowser.util.joinURL(runningPath, pathElement);

      breadcrumbs.push({display: pathElement, path: runningPath});
   });

   return breadcrumbs;
}

// Encode a path for use in a hash.
// We could just do a full encodeURIComponent(), but we can handle leaving
// slashes and spaces alone. This increases readability of the URL.
filebrowser.nav.encodeForHash = function(path) {
   var encodePath = encodeURIComponent(path);

   // Unreplace the slash (%2F) and space (%20).
   return encodePath.replace(/%2F/g, '/').replace(/%20/g, ' ');
}

// Remove the leading hash and decode the path
filebrowser.nav.cleanHashPath = function() {
   return decodeURIComponent(window.location.hash.replace(/^#/, ''));
}
