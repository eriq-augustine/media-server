"use strict";

var filebrowser = filebrowser || {};
filebrowser.filetypes = filebrowser.filetypes || {};

filebrowser.filetypes.templates = filebrowser.filetypes.templates || {};
filebrowser.filetypes.overrides = filebrowser.filetypes.overrides || {};

// TODO(eriq): Icons here.
filebrowser.filetypes.fileClasses = filebrowser.filetypes.fileClasses || {
   'text': {renderFunction: _renderGeneralIFrame},
   'audio': {renderFunction: _renderAudio},
   'image': {renderFunction: _renderImage},
   'general': {renderFunction: _renderGeneralIFrame},
   'html': {renderFunction: _renderGeneralIFrame},
   'video': {renderFunction: _renderVideo},
   'code': {renderFunction: _renderGeneralIFrame},
};

filebrowser.filetypes.extensions = filebrowser.filetypes.extensions || {
   '':      {fileClass: 'text', mime: 'text/plain'}, // Treat no extension files as text.
   'txt':   {fileClass: 'text', mime: 'text/plain'},
   'nfo':   {fileClass: 'text', mime: 'text/plain'},

   'mp3':   {fileClass: 'audio', mime: 'audio/mpeg'},
   'ogg':   {fileClass: 'audio', mime: 'audio/ogg'},

   'jpg':   {fileClass: 'image', mime: 'image/jpeg'},
   'jpeg':  {fileClass: 'image', mime: 'image/jpeg'},
   'png':   {fileClass: 'image', mime: 'image/png'},
   'gif':   {fileClass: 'image', mime: 'image/gif'},
   'tiff':  {fileClass: 'image', mime: 'image/tiff'},
   'svg':   {fileClass: 'image', mime: 'image/svg+xml'},

   'pdf':   {fileClass: 'general', mime: 'application/pdf'},

   'html':  {fileClass: 'html', mime: 'text/html'},

   'mp4':   {fileClass: 'video', mime: 'video/mp4'},
   'm4v':   {fileClass: 'video', mime: 'video/mp4'},
   'ogv':   {fileClass: 'video', mime: 'video/ogg'},
   'ogx':   {fileClass: 'video', mime: 'video/ogg'},
   'webm':  {fileClass: 'video', mime: 'video/webm'},
   'avi':   {fileClass: 'video', mime: 'video/mp4'},
   'flv':   {fileClass: 'video', mime: 'video/mp4'},
   'mkv':   {fileClass: 'video', mime: 'video/mp4'},

   'sh':    {fileClass: 'code', mime: 'application/x-sh'},
   'java':  {fileClass: 'code', mime: 'text/x-java-source'},
   'rb':    {fileClass: 'code', mime: 'text/x-script.ruby'},
   'py':    {fileClass: 'code', mime: 'text/x-script.phyton'},
};

filebrowser.filetypes.renderHTML = function(file) {
   // TODO(eriq): More error
   if (file.isDir) {
      console.log("Error: Expecting a file, got a directory.");
      return {html: ""};
   }

   var fileClass = undefined;
   if (filebrowser.filetypes.extensions[file.extension]) {
      fileClass = filebrowser.filetypes.extensions[file.extension].fileClass;
   }

   if (!filebrowser.filetypes.fileClasses[fileClass]) {
      // TODO(eriq): More error
      console.log("Error: Unknown extension: " + file.extension);
      return {html: _renderGeneralIFrame(file)};
   }

   var renderInfo = filebrowser.filetypes.fileClasses[fileClass].renderFunction(file);
   if (typeof renderInfo === 'string') {
      return {html: renderInfo};
   }

   return renderInfo;
}

// Rendering functions can return a string (the html to be rendered) or
// an object {html: '', callback: ()}.
filebrowser.filetypes.registerRenderOverride = function(fileClass, renderFunction) {
   if (!filebrowser.filetypes.fileClasses[fileClass]) {
      // TODO(eriq): Better logging
      console.log("Cannot register override, unknown fileClass: " + fileClass);
      return false;
   }

   filebrowser.filetypes.fileClasses[fileClass].renderFunction = renderFunction;

   return true;
}

function _renderGeneralIFrame(file) {
   return filebrowser.filetypes.templates.generalIFrame
      .replace('{{RAW_URL}}', file.directLink);
}

function _renderAudio(file) {
   return filebrowser.filetypes.templates.audio
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{MIME}}', filebrowser.filetypes.extensions[file.extension].mime);
}

function _renderVideo(file) {
   return filebrowser.filetypes.templates.video
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{MIME}}', filebrowser.filetypes.extensions[file.extension].mime);
}

function _renderImage(file) {
   return filebrowser.filetypes.templates.image
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{BASE_NAME}}', file.basename);
}

function _renderUnsupported(file) {
   return filebrowser.filetypes.templates.unsupported.replace('{{EXTENSION}}', file.extension);
}

filebrowser.filetypes.templates.generalIFrame = `
   <iframe src='{{RAW_URL}}'>
      Browser Not Supported
   </iframe>
`;

filebrowser.filetypes.templates.audio = `
   <audio controls>
      <source src='{{RAW_URL}}' type='{{MIME}}'>
      Browser Not Supported
   </audio>
`;

filebrowser.filetypes.templates.image = `
   <img src='{{RAW_URL}}' title='{{BASE_NAME}}' alt='{{BASE_NAME}}'>
`;

filebrowser.filetypes.templates.unsupported = `
   <h2>File type ({{EXTENSION}}) is not supported.</h2>
`;

filebrowser.filetypes.templates.video = `
   <video controls>
      <source src='{{RAW_URL}}' type='{{MIME}}'>
      Browser Not Supported
   </video>
`;
