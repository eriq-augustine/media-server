"use strict";

var filebrowser = filebrowser || {};
filebrowser.util = filebrowser.util || {};

filebrowser.util.joinURL = function(base, addition) {
   if (!base || base == '/') {
      return '/' + addition;
   }

   return base + '/' + addition;
}
