"use strict";

var filebrowser = filebrowser || {};
filebrowser.initFields = filebrowser.initFields || {};

// Valid options: {cacheValidator: func(cacheListing), renderOverrides: {fileClass: func(file)}}
filebrowser.init = function(containerId, fetchFunction, options) {
   options = options || {};

   filebrowser.customFetch = fetchFunction;

   filebrowser.containerId = containerId;
   filebrowser.tableId = containerId + '-tablesorter';

   filebrowser.containerQuery = '#' + filebrowser.containerId;
   filebrowser.bodyContentQuery = filebrowser.containerQuery + ' .filebrowser-body-content';
   filebrowser.tableQuery = '#' + filebrowser.tableId;
   filebrowser.breadcrumbQuery = filebrowser.containerQuery + ' .filebrowser-breadcrumbs-area';

   filebrowser.initFields._parseOptions(options);
}

filebrowser.initFields._parseOptions = function(options) {
   if (options.hasOwnProperty('renderOverrides')) {
      for (var fileClass in options.renderOverrides) {
         if(options.renderOverrides.hasOwnProperty(fileClass)) {
            filebrowser.filetypes.registerRenderOverride(fileClass, options.renderOverrides[fileClass]);
         }
      }
   }

   if (options.hasOwnProperty('cacheValidator')) {
      filebrowser.cache.customValidation = options.cacheValidator;
   }
}
