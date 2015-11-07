"use strict";

var templates = templates || {};

var extensions = {
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

function getHTML(file) {
   // TODO(eriq): More error
   if (file.isDir) {
      console.log("Error: Expecting a file, got a directory.");
      return "";
   }

   var template = undefined;
   if (extensions[file.extension]) {
      template = extensions[file.extension].template;
   }

   switch (template) {
      case 'text':
         return generalIFrame(file);
      case 'code':
         return generalIFrame(file);
      case 'audio':
         return audio(file);
      case 'image':
         return image(file);
      case 'html':
         return generalIFrame(file);
      case 'video':
         return video(file);
      case 'general':
         return generalIFrame(file);
      case 'unsupported':
         return unsupported(file);
      default:
         // TODO(eriq): More error
         console.log("Error: Unknown extension: " + file.extension);
         return generalIFrame(file);
   }
}

function generalIFrame(file) {
   return templates.generalIFrame.replace('{{RAW_URL}}', file.directLink);
}

function audio(file) {
   return templates.audio
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{MIME}}', extensions[file.extension].mime);
}

function video(file) {
   return templates.video
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{MIME}}', extensions[file.extension].mime);
}

function image(file) {
   return templates.image
      .replace('{{RAW_URL}}', file.directLink)
      .replace('{{BASE_NAME}}', file.basename);
}

function unsupported(file) {
   return templates.unsupported.replace('{{EXTENSION}}', file.extension);
}

templates.generalIFrame = `
   <iframe src='{{RAW_URL}}'>
      Browser Not Supported
   </iframe>
`;

templates.audio = `
   <audio controls>
      <source src='{{RAW_URL}}' type='{{MIME}}'>
      Browser Not Supported
   </audio>
`;

templates.image = `
   <img src='{{RAW_URL}}' title='{{BASE_NAME}}' alt='{{BASE_NAME}}'>
`;

templates.unsupported = `
   <h2>File type ({{EXTENSION}}) is not supported.</h2>
`;

templates.video = `
   <video controls>
      <source src='{{RAW_URL}}' type='{{MIME}}'>
      Browser Not Supported
   </video>
`;
