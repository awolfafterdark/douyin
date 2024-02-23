// @Title Connection interface
// @Description tcp websocket
// @Author Wangwengang 2021/8/17 4:22 pm
// @Update Wangwengang 2021/8/17 4:22 pm
package fay

import (
"context"
"io"
"net"
)

type Connection interface {
Start() //Start the connection and let the current connection start working
Stop() //Stop the connection and end the current connection state M
Context() context.Context //Return ctx, used for user-defined go process to obtain connection exit status

Write(data []byte) error
GetReader() (r io.Reader, err error)
ConnClose() // Close the socket connection
RemoteAddr() net.Addr //Get remote client address information

StartWriter()
StartReader()

//GetTcpConnection() *net.TCPConn //Get the original socket TCPConn from the current connection
//GetWsConnection() *websocket.Conn // Get the original websocket Conn from the current connection
GetConnID() uint64 //Get the current connection ID

SendMsg(msgID uint32, data []byte) error //Send Message data directly to the remote TCP client (no buffering)
SendBuffMsg(data []byte) error //Send Message data directly to the remote TCP client (with buffering)

SetProperty(key string, value interface{}) //Set link properties
GetProperty(key string) (interface{}, error) //Get link properties
RemoveProperty(key string) //Remove link properties
}
