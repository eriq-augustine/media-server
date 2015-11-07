"use strict";

var filebrowser = filebrowser || {};
filebrowser.filetypes = filebrowser.filetypes || {};

filebrowser.filetypes.templates = filebrowser.filetypes.templates || {};

filebrowser.filetypes.extensions = filebrowser.filetypes.extensions || {
   '':      {template: 'text', mime: 'text/plain'}, // Treat no extension files as text.
   'txt':   {template: 'text', mime: 'text/plain'},
   'nfo':   {template: 'text', mime: 'text/plain'},

   'mp3':   {template: 'audio', mime: 'audio/mpeg'},
   'ogg':   {template: 'audio', mime: 'audio/ogg'},

   'jpg':   {template: 'image', mime: 'image/jpeg'},
   'jpeg':  {template: 'image', mime: 'image/jpeg'},
   'png':   {template: 'image', mime: 'image/png'},
   'gif':   {template: 'image', mime: 'image/gif'},
   'tiff':  {template: 'image', mime: 'image/tiff'},
   'svg':   {template: 'image', mime: 'image/svg+xml'},

   'pdf':   {template: 'general', mime: 'application/pdf'},

   'html':  {template: 'html', mime: 'text/html'},

   'mp4':   {template: 'video', mime: 'video/mp4'},
   'm4v':   {template: 'video', mime: 'video/mp4'},
   'ogv':   {template: 'video', mime: 'video/ogg'},
   'ogx':   {template: 'video', mime: 'video/ogg'},
   'webm':  {template: 'video', mime: 'video/webm'},
   'avi':   {template: 'video', mime: 'video/mp4'},
   'flv':   {template: 'video', mime: 'video/mp4'},
   'mkv':   {template: 'video', mime: 'video/mp4'},

   'sh':    {template: 'code', mime: 'application/x-sh'},
   'java':  {template: 'code', mime: 'text/x-java-source'},
   'rb':    {template: 'code', mime: 'text/x-script.ruby'},
   'py':    {template: 'code', mime: 'text/x-script.phyton'},
};

filebrowser.filetypes.renderHTML = function(file) {
   // TODO(eriq): More error
   if (file.isDir) {
      console.log("Error: Expecting a file, got a directory.");
      return "";
   }

   var template = undefined;
   if (filebrowser.filetypes.extensions[file.extension]) {
      template = filebrowser.filetypes.extensions[file.extension].template;
   }

   switch (template) {
      case 'text':
         return filebrowser.filetypes._generalIFrame(file);
      case 'code':
         return filebrowser.filetypes._generalIFrame(file);
      case 'audio':
         return filebrowser.filetypes._audio(file);
      case 'image':
         return filebrowser.filetypes._image(file);
      case 'html':
         return filebrowser.filetypes._generalIFrame(file);
      case 'video':
         return filebrowser.filetypes._video(file);
      case 'general':
         return filebrowser.filetypes._generalIFrame(file);
      case 'unsupported':
         return filebrowser.filetypes._unsupported(file);
      default:
         // TODO(eriq): More error
         console.log("Error: Unknown extension: " + file.extension);
         return filebrowser.filetypes._generalIFrame(file);
   }
}

filebrowser.filetypes._generalIFrame = function(file) {
   return filebrowser.filetypes.templates.generalIFrame
      .replace('{{RAW_URL}}', file.directLink);
}

filebrowser.filetypes._audio = function(file) {
   return filebrowser.filetypes.templates.audio
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{MIME}}', filebrowser.filetypes.extensions[file.extension].mime);
}

filebrowser.filetypes._video = function(file) {
   return filebrowser.filetypes.templates.video
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{MIME}}', filebrowser.filetypes.extensions[file.extension].mime);
}

filebrowser.filetypes._image = function(file) {
   return filebrowser.filetypes.templates.image
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{BASE_NAME}}', file.basename);
}

filebrowser.filetypes._unsupported = function(file) {
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
