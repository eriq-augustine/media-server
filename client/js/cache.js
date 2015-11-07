"use strict";

// Start empty.
var fileCache = fileCache || {};

function cacheFind(path) {
   return cacheFindHelper(path.replace(/\/$/, '').split('/'), fileCache);
}

function cacheFindHelper(pathArray, cache) {
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
   return cacheFindHelper(pathArray, cache[element].children);
}

function listingFromCache(path) {
   var cachedListing = cacheFind(path);

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

   return cachedListing;
}

// "Adding" a file to the cache is usually just refreshing or enhancing it's metadata.
function addFileToCache(path, file) {
   fileCache = addFileToCacheHelper(path.replace(/\/$/, '').split('/'), fileCache, file);

   // TEST
   console.log('--- Post Add File To Cache');
   console.log(fileCache);
}

// |cache| is an object of children.
function addFileToCacheHelper(pathArray, cache, file) {
   var pathPart = pathArray.shift();

   // Found the proper place.
   // Put this file into the cahce.
   if (pathArray.length == 0) {
      // TODO(eriq): Fill in this filed as sson as we convert the file.
      file.detailsCached = true;
      cache[pathPart] = file;
      return cache;
   }

   // If we can't find the next part of the path but are not at the proper place,
   // then this means that it has not been cached.
   // Make a false entry for the directory and continue on.
   if (!cache.hasOwnProperty(pathPart)) {
      cache[pathPart] = new Dir(pathPart);
   }

   // Move on to the next dir.
   cache[pathPart].children = addFileToCacheHelper(pathArray, cache[pathPart].children, file);
   return cache;
}

function addDirToCache(path, files) {
   fileCache = addDirToCacheHelper(path.replace(/\/$/, '').split('/'), fileCache, files);

   // TEST
   console.log('--- Post Add Dir To Cache');
   console.log(fileCache);
}

// |cache| is an object of children.
function addDirToCacheHelper(pathArray, cache, files) {
   // TEST
   console.log(cache);

   var pathPart = pathArray.shift();

   // If we can't find the next part of the path, then this means that it has not been cached.
   // Make a false entry for the directory and continue on.
   if (!cache.hasOwnProperty(pathPart)) {
      cache[pathPart] = new Dir(pathPart);
   }

   // Found the proper place.
   // Put all the kids in the cache and mark it as having cached kids.
   if (pathArray.length == 0) {
      cache[pathPart].detailsCached = true;

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
   cache[pathPart].children = addDirToCacheHelper(pathArray, cache[pathPart].children, files);
   return cache;
}

function loadCache(path, callback) {
   fetch(path, function(isDir, data) {
      if (isDir) {
         addDirToCache(path, data);
      } else {
         addFileToCache(path, data);
      }

      callback();
   });
}
