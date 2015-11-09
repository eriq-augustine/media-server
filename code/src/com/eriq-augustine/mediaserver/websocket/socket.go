package websocket;

import (
   "encoding/json"
   "time"

   ws "golang.org/x/net/websocket"

   "com/eriq-augustine/mediaserver/cache"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/messages"
   "com/eriq-augustine/mediaserver/util"
)

const (
   MESSAGE_SIZE = 2048
   REFRESH_DURATION_SEC = 5 // Time between updates to the client.
)

// Keep track of all incoming connections.
var connectionId int;
var connections map[int]*WebSocketInfo;

type WebSocketInfo struct {
   Socket *ws.Conn
   LastMessageHash string // We will keep track of the last message so we don't send repeats.
}

func init() {
   connectionId = 0;
   connections = make(map[int]*WebSocketInfo);

   go sendUpdates();
}

// The entrypoint for the router.
func SocketHandler(socket *ws.Conn) {
   var id = connectionId;
   connectionId++;
   connections[id] = &WebSocketInfo{socket, ""};

   // The client's encoding information will get initialized in the next update cycle.

   var rawMsg []byte = make([]byte, MESSAGE_SIZE);
   for {
      size, err := socket.Read(rawMsg);
      if (err != nil) {
         log.ErrorE("Unable to read websocket messgae", err);
         continue;
      }

      var msg map[string]interface{};
      err = json.Unmarshal(rawMsg[0:size], &msg);
      if (err != nil) {
         log.ErrorE("Unable to unmarshal websocket message", err);
         continue;
      }

      // Right now, we are not actually expecting any new messages.
   }
}

func sendUpdates() {
   var msg *messages.CacheStatus;
   var msgJSON string;
   var msgHash string;

   for {
      time.Sleep(REFRESH_DURATION_SEC * time.Second);

      if (len(connections) == 0) {
         continue;
      }

      msg = getCacheStatus();
      msgJSON, _ = util.ToJSON(msg);
      msgHash = util.SHA1Hex(msgJSON);

      for _, socketInfo := range(connections) {
         if (socketInfo.LastMessageHash == msgHash) {
            continue;
         }

         socketInfo.LastMessageHash = msgHash;

         _, err := socketInfo.Socket.Write([]byte(msgJSON));
         if (err != nil) {
            log.ErrorE("Error sending message: " + msgJSON, err);
            continue;
         }
      }
   }
}

func getCacheStatus() *messages.CacheStatus {
   return messages.NewCacheStatus(cache.GetProgress(), cache.GetQueue(), cache.GetRecentEncodes(-1));
}
