package main

import (
"bytes"
"encoding/json"
"errors"
"github.com/wwengg/douyin/fay/impl"
"github.com/wwengg/douyin/model"
"github.com/wwengg/douyin/proto"
"github.com/wwengg/douyin/utils"
"io"
"log"
"net/http"
"strings"

"github.com/elazarl/goproxy"
)

func main() {
utils.ConfigureCA()
proxy := goproxy.NewProxyHttpServer()
proxy.Verbose = false

fayProxyServer := impl.NewFayProxyServer()
go fayProxyServer.StartWebsocket()

//ws data processing
proxy.AddWebsocketHandler(func(data []byte, direction goproxy.WebsocketDirection, ctx *goproxy.ProxyCtx) (reply []byte) {
reply=data
if len(data) == 0 {
return
}
if data[0] != 0x08 {
return
}
wssResponse := proto.WssResponse{}
if err := wssResponse.XXX_Unmarshal(data); err == nil {
//Detect packet format
if v, ok := wssResponse.Headers["compress_type"]; !ok && v != "gzip" {
return
}
//Decompress gzip
deData, err := utils.GzipDecode(wssResponse.Payload)
if err != nil {
ctx.Logf("gzip decompression failed")
return
}
res := proto.Response{}
if err = res.XXX_Unmarshal(deData); err != nil {
return
}
for _, message := range res.Messages {
fayProxyServer.DoMessage(message)
}
}
return
})

proxy.WebSocketHandler = func(dst io.Writer, src io.Reader, direction goproxy.WebsocketDirection, ctx *goproxy.ProxyCtx) error {
fullPacket := make([]byte, 0)
buf := make([]byte, 32*1024)
var err error = nil
for {
nr, er := src.Read(buf)

if er != nil {
if er != io.EOF {
err = err
}
break
}

if nr > 0 {
fullPacket = append(fullPacket, buf[:nr]...)
websocketPacket := model.NewWebsocketPacket(fullPacket)

if !websocketPacket.Valid {
continue
}

websocketPacket.Payload = proxy.FilterWebsocketPacket(websocketPacket.Payload, direction, ctx)
encodedPacket := websocketPacket.Encode()
nw, ew := dst.Write(encodedPacket)
fullPacket = fullPacket[websocketPacket.PacketSize:]

if nw < 0 || len(encodedPacket) < nw {
nw = 0
if ew == nil {
ew = errors.New("invalid write result")
}
}
if ew != nil {
err = ew
break
}
if len(encodedPacket) != nw {
err = io.ErrShortWrite
break
}
}
}
return err
}

proxy.OnRequest(goproxy.ReqHostIs("webcast.amemv.com:443", "frontier-im.douyin.com:443", "webcast100-ws-web-lq.amemv.com:443", "webcast3-ws -web-lf.douyin.com:443", "webcast3-ws-web-hl.douyin.com:443")).
HandleConnect(goproxy.AlwaysMitm)

proxy.OnResponse(goproxy.UrlHasPrefix("httpswebcast.amemv.com:443/webcast/room/create/")).DoFunc(
func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
buf, _ := io.ReadAll(resp.Body)
responseStream := io.NopCloser(bytes.NewBuffer(buf))
rtmpLive := RtmpLive{}
json.Unmarshal(buf, &rtmpLive)

log.Println(rtmpLive)
url := rtmpLive.Data.StreamUrl.RtmpPushUrl
array := strings.Split(url, "/")
secret := array[len(array)-1]
serverName := strings.Split(url, secret)[0]
log.Printf(`Server: %s`, serverName)
log.Printf(`Push code: %s`, secret)
resp.Body = responseStream
return resp
},
)
log.Println("The software is ready, please start [Live Broadcast Companion] and click [Start Live Broadcast]")
log.Fatal(http.ListenAndServe(":8001", proxy))
}

typeRtmpLive struct {
Data struct {
StreamUrl struct {
RtmpPushUrl string `json:"rtmp_push_url"`
} `json:"stream_url"`
} `json:"data"`
}
