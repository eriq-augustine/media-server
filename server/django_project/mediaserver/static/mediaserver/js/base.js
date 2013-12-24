'use strict';

window.mediaserver = {};

document.addEventListener('DOMContentLoaded', function() {
   // Hide the right-pane until the websocket fills it.
   $('.right-pane').hide();

   // Auto-focus the first form field.
   $('form:first *:input[type!=hidden]:first').focus();
});
