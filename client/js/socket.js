"use strict";

var mediaserver = mediaserver || {};
mediaserver.socket = mediaserver.socket || {};

mediaserver.socket.init = function(path) {
   mediaserver.socket._socket = new WebSocket(path);

   mediaserver.socket._socket.onmessage = mediaserver.socket._onmessage;
   mediaserver.socket._socket.onopen = mediaserver.socket._onopen;
   mediaserver.socket._socket.onclose = mediaserver.socket._onclose;
   mediaserver.socket._socket.onerror = mediaserver.socket._onerror;
}

mediaserver.socket._onmessage = function(ev) {
   if (!ev.data) {
      return;
   }

   var jsonData = JSON.parse(ev.data);

   if (jsonData && jsonData.Success) {
      mediaserver.renderEncodeActivity(jsonData);
   }
}

mediaserver.socket._onopen = function(ev) {
   // TODO(eriq): Better logging.
   console.log("Socket Opened");
}

mediaserver.socket._onclose = function(ev) {
   console.log("Socket Closed");
}

mediaserver.socket._onerror = function(ev) {
   // TODO(eriq): More logging.
   console.log("Websocket error");
   console.log(ev);
}
