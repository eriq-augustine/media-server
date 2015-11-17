package messages;

const (
   SOCKET_MESSAGE_TYPE_INIT = "init"
)

type SocketInit struct {
   Type string
   Token string
}
