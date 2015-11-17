package websocket;

import (
   "encoding/json"
   "io"
   "strconv"
   "time"

   ws "golang.org/x/net/websocket"

   "com/eriq-augustine/mediaserver/auth"
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
var nextConnectionId int;
var connections map[int]*WebSocketInfo;

type WebSocketInfo struct {
   Socket *ws.Conn
   LastMessageHash string // We will keep track of the last message so we don't send repeats.
}

func init() {
   nextConnectionId = 0;
   connections = make(map[int]*WebSocketInfo);

   go sendUpdates();
}

// The entrypoint for the router.
func SocketHandler(socket *ws.Conn) {
   var id = nextConnectionId;
   nextConnectionId++;

   // The client's encoding information will get initialized in the next update cycle.

   // Defer closing the connection and removing the connection from the pool.
   defer func(id int) {
      socketInfo, exists := connections[id];
      if (exists) {
         delete(connections, id);
         socketInfo.Socket.Close();
      }
   }(id);

   var rawMsg []byte = make([]byte, MESSAGE_SIZE);
   for {
      size, err := socket.Read(rawMsg);
      if (err == io.EOF) {
         // No problem here.
         log.Debug("Client closed websocket: " + strconv.Itoa(id));
         break;
      }

      if (err != nil) {
         log.ErrorE("Unable to read websocket messgae", err);
         break;
      }

      var msg map[string]interface{};
      err = json.Unmarshal(rawMsg[0:size], &msg);
      if (err != nil) {
         log.ErrorE("Unable to unmarshal websocket message", err);
         continue;
      }

      _, exists := msg["Type"];
      if (!exists) {
         log.Error("Socket message does not have a type");
         continue;
      }

      stringType, ok := msg["Type"].(string);
      if (!ok) {
         log.Error("Socket message type is not a string");
         continue;
      }

      if (stringType == messages.SOCKET_MESSAGE_TYPE_INIT) {
         var initMsg messages.SocketInit;
         err = json.Unmarshal(rawMsg[0:size], &initMsg);
         if (err != nil) {
            log.ErrorE("Unable to unmarshal socket init", err);
            continue;
         }

         // Validate the token and add the connection to the pool.
         _, err = auth.ValidateToken(initMsg.Token);
         if (err == nil) {
            // Don't worry about double adding.
            connections[id] = &WebSocketInfo{socket, ""};
         }
      } else {
         log.Error("Unknown socket message type: " + stringType);
         continue;
      }
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
