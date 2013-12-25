'use strict';

document.addEventListener('DOMContentLoaded', function() {
   window.mediaserver.socket = new Socket();
});

// TODO(eriq): Get port from settings?
Socket.SERVER = 'ws://' + document.domain + ':6060';
// Socket.VIEW_BASE = 'http://' + document.domain + ':5050/view'
Socket.VIEW_BASE = 'http://' + document.domain + '/view'

function Socket() {
   this.ws = new WebSocket(Socket.SERVER);

   this.ws.onmessage = this.onMessage.bind(this);
   this.ws.onclose = this.onClose.bind(this);
   this.ws.onopen = this.onOpen.bind(this);
   this.ws.onerror = this.onError.bind(this);
}

Socket.prototype.onMessage = function(messageEvent) {
   var message = null;
   try {
      message = JSON.parse(messageEvent.data);
   } catch (ex) {
      error('Server message does not parse.');
      return false;
   }

   switch (message.type) {
      case 'ENCODE_UPDATE':
         this.update_encode_info(message['info']);
         break;
      default:
         console.log('ERROR: Unknown Message Type: ' + message.type);
         break;
   }

   return true;
};

Socket.prototype.onClose = function(messageEvent) {
   console.log("Connection to server closed.");
   return true;
};

Socket.prototype.onOpen = function(chosenPattern, messageEvent) {
   console.log("Connection to server opened.");
   return true;
};

Socket.prototype.onError = function(messageEvent) {
   console.log("WS Error: " + JSON.stringify(messageEvent));
   return true;
};

Socket.prototype.update_encode_info = function(info) {
   var html = [];
   var i;
  
   if (info['encode_queue'].length == 0 && info['recent_cache'].length == 0) {
      $('.right-pane').hide();
      $('.page-content').removeClass('page-content-narrow');
      return;
   } else {
      $('.right-pane').show();
      $('.page-content').addClass('page-content-narrow');
   }

   if (info['encode_queue'].length > 0) {
      html.push("<div class='encode-queue'>");
      html.push("   <p>Encode Queue</p>");

      for (var i in info['encode_queue']) {
         var name = info['encode_queue'][i]['name'];
         var time = info['encode_queue'][i]['time'];
         var url = Socket.VIEW_BASE + '/' + info['encode_queue'][i]['path'];

         var progress = "<p class='encode-progress'>Queued</p>";
         if ('progress' in info['encode_queue'][i] &&
             info['encode_queue'][i]['progress']['total'] > 0) {
            progress = "<p class='encode-progress'>Encoding... <span class='red-text'>" +
                        info['encode_queue'][i]['progress']['current'] +
                        "</span><span>s / </span><span class='green-text'>" +
                        info['encode_queue'][i]['progress']['total'] +
                        "</span>s</p>";
         }

         html.push("   <div class='encode-item'>");
         html.push("      <div class='encode-item-name'>");
         html.push("         <a href='" + url + "'>" + name + "</a>");
         html.push("      </div>");

         if (progress) {
            html.push("      " + progress);
         }

         html.push("      <div class='encode-item-time'>");
         html.push("         " + time);
         html.push("      </div>");
         html.push("   </div>");
      }

      html.push("</div>");

      if (info['recent_cache'].length > 0) {
         html.push("   <hr />");
      }
   }

   if (info['recent_cache'].length > 0) {
      html.push("<div class='recent-cache'>");
      html.push("   <p>Recently Encoded</p>");

      for (var i in info['recent_cache']) {
         var name = info['recent_cache'][i]['name'];
         var time = info['recent_cache'][i]['time'];
         var url = Socket.VIEW_BASE + '/' + info['recent_cache'][i]['path'];

         html.push("   <div class='encode-item'>");
         html.push("      <div class='encode-item-name'>");
         html.push("         <a href='" + url + "'>" + name + "</a>");
         html.push("      </div>");
         html.push("      <div class='encode-item-time'>");
         html.push("         " + time);
         html.push("      </div>");
         html.push("   </div>");
      }

      html.push("</div>");
   }

   $('.right-pane').html(html.join("\n"));
};
