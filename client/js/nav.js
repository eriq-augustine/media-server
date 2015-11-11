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

   // TODO(eriq): Whether or not to change the title (and hash) should be an option.
   // Change the page's title.
   document.title = filebrowser.util.basename(path);

   // Update the breadcrumbs.
   filebrowser.view.loadBreadcrumbs(filebrowser.nav._buildBreadcrumbs(path));
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

// Remove the leading hash and decode the path
filebrowser.nav.cleanHashPath = function() {
   return decodeURIComponent(window.location.hash.replace(/^#/, ''));
}
