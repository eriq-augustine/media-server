"use strict";

var filebrowser = filebrowser || {};
filebrowser.cache = filebrowser.cache || {};

// Start empty.
filebrowser.cache._fileCache = filebrowser.cache._fileCache || {};

filebrowser.cache.listingFromCache = function(path) {
   var cachedListing = filebrowser.cache._cacheFind(path);

   // If there was nothing in the cache, report a miss.
   if (!cachedListing) {
      return undefined;
   }

   // If this listing has not been fully cached (like a skeleton directory
   // created when adding a child to the cache or a file that only has the
   // information from a dirent, then report as a miss so we fetch the full data.
   if (!cachedListing.detailsCached) {
      return undefined;
   }

   // Check for any cutom invalidation.
   if (!filebrowser.cache.customValidation(cachedListing)) {
      return undefined;
   }

   return cachedListing;
}

// Users can override this to manually invalidate a cache intry.
// Return true if the entry is valid.
filebrowser.cache.customValidation = function(cachedListing) {
   return true;
}

filebrowser.cache.loadCache = function(path, callback) {
   filebrowser.customFetch(path, function(isDir, data) {
      if (isDir) {
         filebrowser.cache._addDirToCache(path, data);
      } else {
         filebrowser.cache._addFileToCache(path, data);
      }

      callback();
   });
}

filebrowser.cache._cacheFind = function(path) {
   return filebrowser.cache._cacheFindHelper(path.replace(/\/$/, '').split('/'), filebrowser.cache._fileCache);
}

filebrowser.cache._cacheFindHelper = function(pathArray, cache) {
   // Ran out of places to look.
   if (!pathArray || pathArray.length == 0 || !cache || Object.keys(cache).length == 0) {
      return undefined;
   }

   // Not here.
   // If this is a legitimite dirent, then this should only happen if the parent has not been cached.
   if (!cache.hasOwnProperty(pathArray[0])) {
      return undefined;
   }

   // If there is nothing else left in the path, then we found it.
   if (pathArray.length == 1) {
      return cache[pathArray[0]];
   }

   // Keep looking in the next kid.
   var element = pathArray.shift();
   return filebrowser.cache._cacheFindHelper(pathArray, cache[element].children);
}

// "Adding" a file to the cache is usually just refreshing or enhancing it's metadata.
filebrowser.cache._addFileToCache = function(path, file) {
   filebrowser.cache._fileCache = filebrowser.cache._addFileToCacheHelper(path.replace(/\/$/, '').split('/'), filebrowser.cache._fileCache, file);
}

// |cache| is an object of children.
filebrowser.cache._addFileToCacheHelper = function(pathArray, cache, file) {
   var pathPart = pathArray.shift();

   // Found the proper place.
   // Put this file into the cahce.
   if (pathArray.length == 0) {
      // TODO(eriq): Fill in this filed as sson as we convert the file.
      file.detailsCached = true;
      file.cacheTime = new Date();
      cache[pathPart] = file;
      return cache;
   }

   // If we can't find the next part of the path but are not at the proper place,
   // then this means that it has not been cached.
   // Make a false entry for the directory and continue on.
   if (!cache.hasOwnProperty(pathPart)) {
      cache[pathPart] = new filebrowser.Dir(pathPart);
   }

   // Move on to the next dir.
   cache[pathPart].children = filebrowser.cache._addFileToCacheHelper(pathArray, cache[pathPart].children, file);
   return cache;
}

filebrowser.cache._addDirToCache = function(path, files) {
   filebrowser.cache._fileCache = filebrowser.cache._addDirToCacheHelper(path.replace(/\/$/, '').split('/'), filebrowser.cache._fileCache, files);
}

// |cache| is an object of children.
filebrowser.cache._addDirToCacheHelper = function(pathArray, cache, files) {
   var pathPart = pathArray.shift();

   // If we can't find the next part of the path, then this means that it has not been cached.
   // Make a false entry for the directory and continue on.
   if (!cache.hasOwnProperty(pathPart)) {
      cache[pathPart] = new filebrowser.Dir(pathPart);
   }

   // Found the proper place.
   // Put all the kids in the cache and mark it as having cached kids.
   if (pathArray.length == 0) {
      cache[pathPart].detailsCached = true;
      cache[pathPart].cacheTime = new Date();

      // Build the cache for this subtree and return it.
      // If there is an existing cache here, replace it.
      // If we were asked to refresh a parent, then the children's cache is no longer good.
      var kids = {};
      files.forEach(function(file) {
         kids[file.name] = file;
      });

      cache[pathPart].children = kids;
      return cache;
   }

   // Otherwise, just move on to the next dir.
   cache[pathPart].children = filebrowser.cache._addDirToCacheHelper(pathArray, cache[pathPart].children, files);
   return cache;
}
