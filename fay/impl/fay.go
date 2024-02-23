package impl

import (
"encoding/json"
"github.com/gorilla/websocket"
"github.com/wwengg/douyin/fay"
"github.com/wwengg/douyin/proto"
"log"
"net/http"
"time"
)

type FayProxyServer struct {
// Websocket Addr
WsAddr string
//Used to cache the received message ID and determine whether to receive it repeatedly.
Dictionary map[string][]int64
//The link manager of the current Server
ConnMgr fay.ConnManager

GenID uint64
}

func NewFayProxyServer() *FayProxyServer {
return &FayProxyServer{
WsAddr: "127.0.0.1:8888",
Dictionary: make(map[string][]int64),
ConnMgr: NewConnManager(),
}
}

// GetConnMgr gets the link management
func (s *FayProxyServer) GetConnMgr() fay.ConnManager {
return s.ConnMgr
}

// Start Websocket network service
func (s *FayProxyServer) StartWebsocket() {
//logger.ZapLog.Info("Start Websocket server", zap.String("addr", s.WsAddr))
log.Printf("Start Websocket server addr:%s (used to connect to Fay https://github.com/TheRamU/Fay/tree/fay-sales-edition)", s.WsAddr)
httpServer := &http.Server{
Addr: s.WsAddr,
Handler: &WsHandler{upgrader: websocket.Upgrader{
HandshakeTimeout: 0,
ReadBufferSize: 0,
WriteBufferSize: 0,
WriteBufferPool: nil,
Subprotocols: nil,
Error: nil,
CheckOrigin: func(r *http.Request) bool {
return true
},
EnableCompression: false,
},
sv: s},
ReadTimeout: time.Second * time.Duration(60),
WriteTimeout: time.Second * time.Duration(60),
MaxHeaderBytes: 4096,
}
httpServer.ListenAndServe()
}

typeWsHandler struct {
sv *FayProxyServer
upgrader websocket.Upgrader
}

func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
if r.URL.Path != "/" {
http.Error(w, "Not found", http.StatusNotFound)
return
}
if r.Method != "GET" {
http.Error(w, "Method not allowed", 405)
return
}
conn, err := h.upgrader.Upgrade(w, r, nil)
if err != nil {
log.Fatalf("upgrade error,err:%s", err.Error())
return
}

//3.2 Set the server's maximum connection control. If the maximum connection is exceeded, the new connection will be closed.
if h.sv.ConnMgr.Len() >= 10 {
conn.Close()
}

//3.3 Business method to handle the new connection request. At this time, handler and conn should be bound.
h.sv.GenID++
dealConn := NewConnection(h.sv, h.sv.GenID, NewWsProtocol(conn))

h.sv.GetConnMgr().Add(dealConn)
//3.4 Start the processing business of the current link
go dealConn.Start()

}

func (s *FayProxyServer) DoMessage(message *proto.Message) {
list, ok := s.Dictionary[message.Method]
if !ok {
list = []int64{}
} else {
for _, i := range list {
if i == message.MsgId {
return
}
}
}
if len(list) > 300 {
list = []int64{}
}
log.Printf("method:%s,num:%d", message.Method, len(list))
list = append(list, message.MsgId)
s.Dictionary[message.Method] = list
switch message.Method {
// User enters
case "WebcastMemberMessage":
memberMessage := proto.MemberMessage{}
memberMessage.XXX_Unmarshal(message.Payload)
s.onMemberMessgae(memberMessage)
break
	// User Comments
case "WebcastChatMessage":
chatMessage := proto.ChatMessage{}
chatMessage.XXX_Unmarshal(message.Payload)
s.onChatMessage(chatMessage)
break
default:
break
}

}

func (s *FayProxyServer) send(pack *fay.MsgPack) {
data, _ := json.Marshal(pack)
s.GetConnMgr().SendMsgToAllConn(data)
}

func (s *FayProxyServer) onMemberMessgae(message proto.MemberMessage) {
s.send(fay.CreateMsgPack(message, fay.MsgType_JoinRoom))
log.Printf("%s is here", message.User.Nickname)
}

func (s *FayProxyServer) onChatMessage(message proto.ChatMessage) {
s.send(fay.CreateMsgPack(message, fay.MsgType_DanMu))
log.Printf("%s said: %s", message.User.Nickname, message.Content)
}
