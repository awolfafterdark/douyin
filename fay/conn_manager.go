package fay

type ConnManager interface {
Add(conn Connection) //Add link
Remove(conn Connection) //Delete connection
Get(connID uint64) (Connection, error) //Use ConnID to get the link
Len() int //Get the current connection
ClearConn() //Delete and stop all links
SendMsgToAllConn(data []byte) // Send

}
