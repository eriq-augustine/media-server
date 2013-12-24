'use strict';

document.addEventListener('DOMContentLoaded', function() {
   window.mediaserver.socket = new Socket();
});

// TODO(eriq): Get this from settings?
Socket.SERVER = 'ws://' + document.domain + ':6060';

function Socket() {
   this.ws = new WebSocket(Socket.SERVER);

   this.ws.onmessage = this.onMessage.bind(this);
   this.ws.onclose = this.onClose.bind(this);
   this.ws.onopen = this.onOpen.bind(this);
   this.ws.onerror = this.onError.bind(this);
}

Socket.prototype.onMessage = function(messageEvent) {
   //TEST
   console.log(messageEvent.data);

   var message = null;
   try {
      message = JSON.parse(messageEvent.data);
   } catch (ex) {
      error('Server message does not parse.');
      return false;
   }

   /*
   switch (message.Type) {
      case Message.TYPE_START:
         break;
      case Message.TYPE_NEXT_TURN:
         break;
      default:
         // Note: There are messages that are known, but just not expected from the server.
         error('Unknown Message Type: ' + message.Type);
         break;
   }
   */

   return true;
};

Socket.prototype.onClose = function(messageEvent) {
   //TEST
   console.log("Connection to server closed.");

   return true;
};

Socket.prototype.onOpen = function(chosenPattern, messageEvent) {
   console.log("Connection to server opened.");
   //TEST
   // this.ws.send(createInitMessage(chosenPattern));

   return true;
};

Socket.prototype.onError = function(messageEvent) {
   console.log("WS Error: " + JSON.stringify(messageEvent));

   return true;
};

Socket.prototype.close = function() {
   this.ws.close();

   return true;
};

//TEST
// this.ws.send(createMoveMessage(dropGemLocations, boardHash));
