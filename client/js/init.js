"use strict";

var filebrowser = filebrowser || {};
filebrowser.initFields = filebrowser.initFields || {};

filebrowser.initFields._containerTemplate = `
   <div class='filebrowser-head-area'>
      <div class='filebrowser-breadcrumbs-area'>
      </div>
   </div>
   <div class='filebrowser-body-area'>
      <div class='filebrowser-body-content'>
      </div>
   </div>
`

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
   filebrowser.initFields._initHTML(options);
   filebrowser.initFields._initTablesorter();
}

filebrowser.initFields._initHTML = function() {
   $(filebrowser.containerQuery).addClass('filebrowser-container').html(filebrowser.initFields._containerTemplate);
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

filebrowser.initFields._initTablesorter = function() {
   $.tablesorter.addParser({
      id: 'fileSize',
      is: function(s) { // return false so this parser is not auto detected
         return false;
      },
      format: function(data) {
         // Convert the data to bytes for sorting.
         return filebrowser.util.humanToBytes(data);
      },
      type: 'numeric'
   });
}
