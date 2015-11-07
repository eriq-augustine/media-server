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
