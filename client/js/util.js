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

// Get the last component of |path|.
// e.g. the basename of '/a/b/c.txt' is 'c.txt'.
filebrowser.util.basename = function(path) {
   return path.split(/[\\/]/).pop();
}

// See http://stackoverflow.com/questions/190852/how-can-i-get-file-extensions-with-javascript/1203361#1203361
filebrowser.util.ext = function(path) {
   return path.substr((~-path.lastIndexOf(".") >>> 0) + 2);
}

filebrowser.util.formatDate = function(date) {
   date = new Date(date);

   var month = filebrowser.util.zeroPad(date.getMonth() + 1);
   var day = filebrowser.util.zeroPad(date.getDate());
   var hours = filebrowser.util.zeroPad(date.getHours());
   var minutes = filebrowser.util.zeroPad(date.getMinutes());
   var seconds = filebrowser.util.zeroPad(date.getSeconds());

   return '' + date.getFullYear() + '-' + month + '-' + day + ' ' + hours + ':' + minutes + ':' + seconds;
}

filebrowser.util.zeroPad = function(str, desiredLength) {
   desiredLength = desiredLength || 2;

   return ('00000000000000000000000000' + str).slice(desiredLength * -1);
}
