"use strict";

var mediaserver = mediaserver || {};
mediaserver.store = mediaserver.store || {};

mediaserver.store._backend = window.localStorage;

mediaserver.store.TOKEN_KEY = 'api-token';

mediaserver.store.set = function(key, val) {
   if ((typeof val) == 'string') {
      mediaserver.store._backend[key] = val;
   } else {
      mediaserver.store._backend[key] = JSON.stringify(val);
   }
}

mediaserver.store.has = function(key) {
   return mediaserver.store._backend.hasOwnProperty(key);
}

mediaserver.store.get = function(key, defaultValue) {
   if (!mediaserver.store.has(key)) {
      return defaultValue;
   }

   return mediaserver.store._backend[key];
}

mediaserver.store.getObject = function(key, defaultValue) {
   if (!mediaserver.store.has(key)) {
      return defaultValue;
   }

   return JSON.parse(mediaserver.store._backend[key]);
}

mediaserver.store.unset = function(key) {
   if (mediaserver.store.has(key)) {
      delete mediaserver.store._backend[key];
   }
}
