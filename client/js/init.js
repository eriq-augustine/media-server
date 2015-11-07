"use strict";

var filebrowser = filebrowser || {};
filebrowser.init = filebrowser.init || {};

filebrowser.init = function(containerId, fetchFunction) {
   filebrowser.customFetch = fetchFunction;

   filebrowser.containerId = containerId;
   filebrowser.tableId = containerId + '-tablesorter';

   filebrowser.containerQuery = '#' + filebrowser.containerId;
   filebrowser.bodyContentQuery = filebrowser.containerQuery + ' .filebrowser-body-content';
   filebrowser.tableQuery = '#' + filebrowser.tableId;
   filebrowser.breadcrumbQuery = filebrowser.containerQuery + ' .filebrowser-breadcrumbs-area';
}
