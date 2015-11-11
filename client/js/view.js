"use strict";

var filebrowser = filebrowser || {};
filebrowser.view = filebrowser.view || {};

filebrowser.view_arrayToTableRow = function(data, isHeader) {
   isHeader = isHeader | false;
   var cellType = isHeader ? 'th' : 'td';

   var tr = document.createElement('tr');

   data.forEach(function(dataObject) {
      var td = document.createElement(cellType);
      td.appendChild(document.createTextNode(dataObject));
      tr.appendChild(td);
   });

   return tr;
}

filebrowser.view_fileToTableRow = function(file) {
   var typeName = filebrowser.filetypes.getFileClass(file) || 'unknown';
   var data = [file.name, filebrowser.util.formatDate(file.modDate), typeName, filebrowser.util.bytesToHuman(file.size)];
   return filebrowser.view_arrayToTableRow(data, false);
}

filebrowser.view_filesToTable = function(path, files) {
   var table = document.createElement('table');

   var tableHead = document.createElement('thead');
   var headerData = ['Name', 'Date', 'Type', 'Size'];
   tableHead.appendChild(filebrowser.view_arrayToTableRow(headerData, true));
   table.appendChild(tableHead);

   var tableBody = document.createElement('tbody');
   files.forEach(function(file) {
      var row = filebrowser.view_fileToTableRow(file);
      var url = filebrowser.util.joinURL(path, file.name);
      row.setAttribute('data-path', url);
      row.addEventListener('click', filebrowser.nav.changeTarget.bind(window, url));
      tableBody.appendChild(row);
   });
   table.appendChild(tableBody);

   return table;
}

filebrowser.view.clearContent = function() {
   $(filebrowser.bodyContentQuery).empty();
}

filebrowser.view.loadViewer = function(file, path) {
   // TODO(eriq): Re-architect the html some, it's not just a table.
   filebrowser.view.clearContent();

   var renderInfo = filebrowser.filetypes.renderHTML(file);

   $(filebrowser.bodyContentQuery).html(renderInfo.html);

   if (renderInfo.callback) {
      renderInfo.callback();
   }
}

filebrowser.view.reloadTable = function(files, path) {
   var table = filebrowser.view_filesToTable(path, files);
   table.id = filebrowser.tableId;
   table.className = 'tablesorter';

   // TODO(eriq): Better ids
   filebrowser.view.clearContent();
   $(filebrowser.bodyContentQuery).append(table);

   $(filebrowser.tableQuery).tablesorter({
      sortList: [[0,0]],
      widgets: ['zebra']
   });
}

// |breadcrumbs| should be [{display: '', path: ''}, ...].
filebrowser.view.loadBreadcrumbs = function(breadcrumbs) {
   var breadcrumbsElement = document.createElement('div');
   breadcrumbsElement.className = 'filebrowser-breadcrumbs';

   breadcrumbs.forEach(function(breadcrumb, index) {
      var breadcrumbElement = document.createElement('div');
      breadcrumbElement.className = 'filebrowser-breadcrumb';

      // Don'r register a handler for the last element (we are already there).
      if (index != breadcrumbs.length - 1) {
         breadcrumbElement.onclick = filebrowser.nav.changeTarget.bind(window, breadcrumb.path);
      }

      var breadcrumbTextElement = document.createElement('span');
      breadcrumbTextElement.textContent = breadcrumb.display;

      breadcrumbElement.appendChild(breadcrumbTextElement);
      breadcrumbsElement.appendChild(breadcrumbElement);

      // Don't put separators after the first or last elements.
      if (index != 0 && index != breadcrumbs.length - 1) {
         var separator = document.createElement('span');
         separator.className = 'filebrowser-breadcrumb-separator';
         separator.textContent = '/';
         breadcrumbsElement.appendChild(separator);
      }
   });

   $(filebrowser.breadcrumbQuery).empty();
   $(filebrowser.breadcrumbQuery).append(breadcrumbsElement);
}
