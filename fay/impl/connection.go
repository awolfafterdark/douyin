// @Title
// @Description
// @Author Wangwengang 2022/3/30 11:50 pm
// @Update Wangwengang 2022/3/30 11:50 pm
package impl

import (
"context"
"errors"
"fmt"
"github.com/wwengg/douyin/fay"
"io"
"sync"
)

type Connection struct {
FayProxyServer *FayProxyServer
//The currently connected socket TCP socket
//Conn *net.TCPConn
//The ID of the current connection can also be called SessionID. The ID is globally unique.
ConnID uint64
//Message management MsgID and the message management module of the corresponding processing method
//MsgHandler anet2.MsgHandle
//Inform the channel that the link has exited/stopped
ctx context.Context
cancel context.CancelFunc
//Unbuffered pipe, used for reading and writing message communication between two goroutines
msgChanchan[]byte
//There is a buffer pipe for reading and writing message communication between two goroutines.
msgBuffChanchan[]byte

sync.RWMutex
//Link properties
property map[string]interface{}
////Protect the lock of the current property
propertyLock sync.Mutex
//Closed status of the current connection
isClosed bool

fay.Protocol
}

// NewConnection method to create a connection
func NewConnection(fayProxyServer *FayProxyServer, connID uint64, protocol fay.Protocol) fay.Connection {
//Initialize Conn property
c := &Connection{
FayProxyServer: fayProxyServer,
ConnID: connID,
isClosed: false,
//MsgHandler: msgHandler,
msgChan: make(chan []byte),
msgBuffChan: make(chan []byte, 1024),
property: nil,
Protocol: protocol,
}

return c
}

// StartWriter writes messages to Goroutine, and the user sends data to the client
func (c *Connection) StartWriter() {
fmt.Println("[Writer Goroutine is running]")
defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

for {
select {
case data := <-c.msgChan:
//There is data to be written to the client
if err := c.Write(data); err != nil {
fmt.Println("Send Data error:, ", err, " Conn Writer exit")
return
}
//fmt.Printf("Send data succ! data = %+v\n", data)
case data, ok := <-c.msgBuffChan:
if ok {
//There is data to be written to the client
if err := c.Write(data); err != nil {
fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
return
}
} else {
fmt.Println("msgBuffChan is Closed")
break
}
case <-c.ctx.Done():
return
}
}
}

// StartReader reads message Goroutine, used to read data from the client
func (c *Connection) StartReader() {
fmt.Println("[Reader Goroutine is running]")
defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
defer c.Stop()

//Create unpacked and unpacked objects
for {
select {
case <-c.ctx.Done():
return
default:
r, err := c.GetReader()
if err != nil {
//logger.ZapLog.Error("GetReader err:", zap.Error(err))
return
}
if _, err := io.ReadAll(r); err != nil {
fmt.Println("read msg head error ", err)
return
}
}
}
}

// Start starts the connection and lets the current connection start working
func (c *Connection) Start() {
c.ctx, c.cancel = context.WithCancel(context.Background())
//1 Start the Goroutine for the user to read data from the client
go c.StartReader()
//2 Enable Goroutine for writing back the client data process
go c.StartWriter()
//Execute the hook method according to the business that needs to be processed when creating a connection passed in by the user
//c.TCPServer.CallOnConnStart(c)
}

// Stop stops the connection and ends the current connection state M
func (c *Connection) Stop() {
//If the user registers the closing callback service of this link, then the call should be displayed at this moment
//c.TCPServer.CallOnConnStop(c)

c.Lock()
defer c.Unlock()

//If the current link is closed
if c.isClosed == true {
return
}

fmt.Println("Conn Stop()...ConnID = ", c.ConnID)

// Close the socket link
c.ConnClose()
//Close Writer
c.cancel()

//Remove the link from the connection manager
c.FayProxyServer.GetConnMgr().Remove(c)

//Close all pipes of this link
close(c.msgBuffChan)
//Set the flag bit
c.isClosed = true

}

// GetConnID gets the current connection ID
func (c *Connection) GetConnID() uint64 {
return c.ConnID
}

// SendMsg directly sends Message data to the remote TCP client
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
c.RLock()
defer c.RUnlock()
if c.isClosed == true {
return errors.New("connection closed when send msg")
}

//write back to client
c.msgChan <- data

return nil
}

// SendBuffMsg BuffMsg occurs
func (c *Connection) SendBuffMsg(data []byte) error {
c.RLock()
defer c.RUnlock()
if c.isClosed == true {
return errors.New("Connection closed when send buff msg")
}

//write back to client
c.msgBuffChan <- data

return nil
}

// SetProperty sets link properties
func (c *Connection) SetProperty(key string, value interface{}) {
c.propertyLock.Lock()
defer c.propertyLock.Unlock()
if c.property == nil {
c.property = make(map[string]interface{})
}

c.property[key] = value
}

// GetProperty gets the link properties
func (c *Connection) GetProperty(key string) (interface{}, error) {
c.propertyLock.Lock()
defer c.propertyLock.Unlock()

if value, ok := c.property[key]; ok {
return value, nil
}

return nil, errors.New("no property found")
}

// RemoveProperty removes link properties
func (c *Connection) RemoveProperty(key string) {
c.propertyLock.Lock()
defer c.propertyLock.Unlock()

delete(c.property, key)
}

//Return ctx, used for user-defined go process to obtain connection exit status
func (c *Connection) Context() context.Context {
return c.ctx
}
