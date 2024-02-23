package fay

import "encoding/json"

type MsgType uint32

// None = 0,
// Barrage message = 1,
// Like message = 2,
// Enter the live broadcast room = 3,
//Follow news = 4,
// gift message = 5,
// Live broadcast room statistics = 6,
// Fan group message = 7,
// Live broadcast room sharing = 8,
//downcast = 9
const (
MsgType_None MsgType = iota // None
MsgType_DanMu // Danmaku information
MsgType_Dianzan // Like information
MsgType_JoinRoom // Enter the live broadcast room
MsgType_Star // Pay attention to information
MsgType_Gift // Gift information
MsgType_Repo // Live broadcast room statistics
MsgType_FansMsg // Fan information
MsgType_Share // Live broadcast room sharing
MsgType_Offline // Download
)

type FansclubType uint32

//None = 0,
//Fan club upgrade = 1,
//Join the fan club = 2

const (
FansclubType_None FansclubType = iota
FansclubType_Upgrade // Upgrade
FansclubType_JoinFans // Join the fan club
)

type ShareType uint32

// unknown = 0,
// WeChat = 1,
// Moments = 2,
// Weibo = 3,
// QQ space = 4,
// QQ = 5,
// TikTok friends = 112
const (
ShareType_Unkonwn ShareType = iota // unknown
ShareType_Wechat // WeChat
ShareType_WechatQuan // Moments
ShareType_Weibo // Weibo
ShareType_QQZone // qq space
ShareType_QQ // qq
ShareType_DouyinFriend ShareType = 112 // Douyin friends
)

type MsgPack struct {
MsgType MsgType `json:"type"`
Data string `json:"data"`
}

func CreateMsgPack(data interface{}, msgType MsgType) *MsgPack {
jsonData, _ := json.Marshal(data)
return &MsgPack{
MsgType: msgType,
Data: string(jsonData),
}
}
