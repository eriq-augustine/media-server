"use strict";

var filebrowser = filebrowser || {};
filebrowser.util = filebrowser.util || {};

filebrowser.util.joinURL = function(base, addition) {
   // Watch for double roots.
   if (!base || base == '/') {
      if (addition == '/') {
         return '/';
      }

      return '/' + addition;
   }

   return base + '/' + addition;
}

// See http://stackoverflow.com/questions/190852/how-can-i-get-file-extensions-with-javascript/1203361#1203361
filebrowser.util.ext = function(path) {
   return path.substr((~-path.lastIndexOf(".") >>> 0) + 2);
}
