package impl

import (
"errors"
"fmt"
"github.com/wwengg/douyin/fay"
"sync"
)

type ConnManager struct {
connections map[uint64]fay.Connection
connLock sync.RWMutex
}

// NewConnManager creates a link manager
func NewConnManager() *ConnManager {
return &ConnManager{
connections: make(map[uint64]fay.Connection),
}
}

//Add add link
func (connMgr *ConnManager) Add(conn fay.Connection) {

connMgr.connLock.Lock()
//Add conn connection to ConnMananger
connMgr.connections[conn.GetConnID()] = conn
connMgr.connLock.Unlock()

fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

// Remove delete connection
func (connMgr *ConnManager) Remove(conn fay.Connection) {

connMgr.connLock.Lock()
//Delete connection information
delete(connMgr.connections, conn.GetConnID())
connMgr.connLock.Unlock()
fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

// Get uses ConnID to get the link
func (connMgr *ConnManager) Get(connID uint64) (fay.Connection, error) {
connMgr.connLock.RLock()
defer connMgr.connLock.RUnlock()

if conn, ok := connMgr.connections[connID]; ok {
return conn, nil
}

return nil, errors.New("connection not found")

}

// Len gets the current connection
func (connMgr *ConnManager) Len() int {
connMgr.connLock.RLock()
length := len(connMgr.connections)
connMgr.connLock.RUnlock()
return length
}

// ClearConn clears and stops all connections
func (connMgr *ConnManager) ClearConn() {
connMgr.connLock.Lock()

//Stop and delete all connection information
for connID, conn := range connMgr.connections {
		//stop
conn.Stop()
		//delete
delete(connMgr.connections, connID)
}
connMgr.connLock.Unlock()
fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

// ClearOneConn uses ConnID to get a link and delete it
func (connMgr *ConnManager) ClearOneConn(connID uint64) {
connMgr.connLock.Lock()
defer connMgr.connLock.Unlock()

connections := connMgr.connections
if conn, ok := connections[connID]; ok {
		//stop
conn.Stop()
		//delete
delete(connections, connID)
fmt.Println("Clear Connections ID: ", connID, "succeed")
return
}

fmt.Println("Clear Connections ID: ", connID, "err")
return
}

// ClearConn clears and stops all connections
func (connMgr *ConnManager) SendMsgToAllConn(data []byte) {
connMgr.connLock.Lock()

//Stop and delete all connection information
for _, conn := range connMgr.connections {
		//send
_ = conn.SendBuffMsg(data)

}
connMgr.connLock.Unlock()
fmt.Println("Send Data to All Connections successfully: conn num = ", connMgr.Len())
}
