'use strict';

document.addEventListener('DOMContentLoaded', function() {
   // Auto-focus the first form field.
   $('form:first *:input[type!=hidden]:first').focus();
});
