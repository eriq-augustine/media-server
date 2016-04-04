"use strict";

var filebrowser = filebrowser || {};

filebrowser.DirEnt = function(name, modDate, size, isDir, detailsCached) {
   this.name = name;
   this.modDate = modDate;
   this.size = size;
   this.isDir = isDir;
   this.detailsCached = detailsCached;
   this.cacheTime = detailsCached ? new Date() : null;
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

   if (name.indexOf('.') > -1) {
      var nameParts = name.match(/^(.*)\.([^\.]*)$/);
      this.basename = nameParts[1];
      this.extension = nameParts[2].toLowerCase();
   } else {
      this.basename = name;
      this.extension = '';
   }
}

filebrowser.File.prototype = Object.create(filebrowser.DirEnt.prototype);
filebrowser.File.prototype.constructor = filebrowser.File;
