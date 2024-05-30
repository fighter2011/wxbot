package uos

import (
	"github.com/yqchilde/wxbot/engine/robot"
	"strings"
)

// 消息类型
const (
	AppMessage = 6
)

// MessageType 以Go惯用形式定义了PC微信所有的官方消息类型。
// 详见 message_test.go
type MessageType int

// AppMessageType 以Go惯用形式定义了PC微信所有的官方App消息类型。
type AppMessageType int

// https://res.wx.qq.com/a/wx_fed/webwx/res/static/js/index_c7d281c.js
// MSGTYPE_TEXT
// MSGTYPE_IMAGE
// MSGTYPE_VOICE
// MSGTYPE_VERIFYMSG
// MSGTYPE_POSSIBLEFRIEND_MSG
// MSGTYPE_SHARECARD
// MSGTYPE_VIDEO
// MSGTYPE_EMOTICON
// MSGTYPE_LOCATION
// MSGTYPE_APP
// MSGTYPE_VOIPMSG
// MSGTYPE_VOIPNOTIFY
// MSGTYPE_VOIPINVITE
// MSGTYPE_MICROVIDEO
// MSGTYPE_SYS
// MSGTYPE_RECALLED

const (
	MsgTypeText           MessageType = 1     // 文本消息
	MsgTypeImage          MessageType = 3     // 图片消息
	MsgTypeVoice          MessageType = 34    // 语音消息
	MsgTypeVerify         MessageType = 37    // 认证消息
	MsgTypePossibleFriend MessageType = 40    // 好友推荐消息
	MsgTypeShareCard      MessageType = 42    // 名片消息
	MsgTypeVideo          MessageType = 43    // 视频消息
	MsgTypeEmoticon       MessageType = 47    // 表情消息
	MsgTypeLocation       MessageType = 48    // 地理位置消息
	MsgTypeApp            MessageType = 49    // APP消息
	MsgTypeVoip           MessageType = 50    // VOIP消息
	MsgTypeVoipNotify     MessageType = 52    // VOIP结束消息
	MsgTypeVoipInvite     MessageType = 53    // VOIP邀请
	MsgTypeMicroVideo     MessageType = 62    // 小视频消息
	MsgTypeSys            MessageType = 10000 // 系统消息
	MsgTypeRecalled       MessageType = 10002 // 消息撤回
)

const (
	AppMsgTypeText                  AppMessageType = 1      // 文本消息
	AppMsgTypeImg                   AppMessageType = 2      // 图片消息
	AppMsgTypeAudio                 AppMessageType = 3      // 语音消息
	AppMsgTypeVideo                 AppMessageType = 4      // 视频消息
	AppMsgTypeUrl                   AppMessageType = 5      // 文章消息
	AppMsgTypeAttach                AppMessageType = 6      // 附件消息
	AppMsgTypeOpen                  AppMessageType = 7      // Open
	AppMsgTypeEmoji                 AppMessageType = 8      // 表情消息
	AppMsgTypeVoiceRemind           AppMessageType = 9      // VoiceRemind
	AppMsgTypeScanGood              AppMessageType = 10     // ScanGood
	AppMsgTypeGood                  AppMessageType = 13     // Good
	AppMsgTypeEmotion               AppMessageType = 15     // Emotion
	AppMsgTypeCardTicket            AppMessageType = 16     // 名片消息
	AppMsgTypeRealtimeShareLocation AppMessageType = 17     // 地理位置消息
	AppMsgTypeTransfers             AppMessageType = 2000   // 转账消息
	AppMsgTypeRedEnvelopes          AppMessageType = 2001   // 红包消息
	AppMsgTypeReaderType            AppMessageType = 100001 //自定义的消息
)

type Message struct {
	IsAt    bool `json:"isAt"`
	AppInfo struct {
		Type  int    `json:"type"`
		AppID string `json:"appID"`
	} `json:"appInfo"`
	AppMsgType           AppMessageType `json:"appMsgType"`
	HasProductId         int            `json:"hasProductId"`
	ImgHeight            int            `json:"imgHeight"`
	ImgStatus            int            `json:"imgStatus"`
	ImgWidth             int            `json:"imgWidth"`
	ForwardFlag          int            `json:"forwardFlag"`
	MsgType              MessageType    `json:"msgType"`
	Status               int            `json:"status"`
	StatusNotifyCode     int            `json:"statusNotifyCode"`
	SubMsgType           int            `json:"subMsgType"`
	VoiceLength          int            `json:"voiceLength"`
	CreateTime           int64          `json:"createTime"`
	NewMsgId             int64          `json:"newMsgId"`
	PlayLength           int64          `json:"playLength"`
	MediaId              string         `json:"mediaId"`
	MsgId                string         `json:"msgId"`
	EncryFileName        string         `json:"encryFileName"`
	FileName             string         `json:"fileName"`
	FileSize             string         `json:"fileSize"`
	Content              string         `json:"content"`
	FromUserName         string         `json:"fromUserName"`
	OriContent           string         `json:"oriContent"`
	StatusNotifyUserName string         `json:"statusNotifyUserName"`
	Ticket               string         `json:"ticket"`
	ToUserName           string         `json:"toUserName"`
	Url                  string         `json:"url"`
	RecommendInfo        RecommendInfo  `json:"recommendInfo"`
	AttachmentUrl        string         `json:"attachmentUrl"`
	//senderUserNameInGroup string                 `json:"senderUserNameInGroup"`
	//item                  map[string]interface{} `json:"item"`
}

// RecommendInfo 一些特殊类型的消息会携带该结构体信息
type RecommendInfo struct {
	OpCode     int
	Scene      int
	Sex        int
	VerifyFlag int
	AttrStatus int64
	QQNum      int64
	Alias      string
	City       string
	Content    string
	NickName   string
	Province   string
	Signature  string
	Ticket     string
	UserName   string
}

func (m *Message) IsText() bool {
	return m.MsgType == MsgTypeText && m.Url == ""
}

func (m *Message) IsLocation() bool {
	return m.MsgType == MsgTypeText && strings.Contains(m.Url, "api.map.qq.com") && strings.Contains(m.Content, "pictype=location")
}

func (m *Message) IsRealtimeLocation() bool {
	return m.IsRealtimeLocationStart() || m.IsRealtimeLocationStop()
}

func (m *Message) IsRealtimeLocationStart() bool {
	return m.MsgType == MsgTypeApp && m.AppMsgType == AppMsgTypeRealtimeShareLocation
}

func (m *Message) IsRealtimeLocationStop() bool {
	return m.MsgType == MsgTypeSys && m.Content == "位置共享已经结束"
}

func (m *Message) IsPicture() bool {
	return m.MsgType == MsgTypeImage
}

// IsEmoticon 是否为表情包消息
func (m *Message) IsEmoticon() bool {
	return m.MsgType == MsgTypeEmoticon
}

func (m *Message) IsVoice() bool {
	return m.MsgType == MsgTypeVoice
}

func (m *Message) IsFriendAdd() bool {
	return m.MsgType == MsgTypeVerify && m.FromUserName == "fmessage"
}

func (m *Message) IsCard() bool {
	return m.MsgType == MsgTypeShareCard
}

func (m *Message) IsVideo() bool {
	return m.MsgType == MsgTypeVideo || m.MsgType == MsgTypeMicroVideo
}

func (m *Message) IsMedia() bool {
	return m.MsgType == MsgTypeApp
}

// IsRecalled 判断是否撤回
func (m *Message) IsRecalled() bool {
	return m.MsgType == MsgTypeRecalled
}

func (m *Message) IsSystem() bool {
	return m.MsgType == MsgTypeSys
}

func (m *Message) IsNotify() bool {
	return m.MsgType == 51 && m.StatusNotifyCode != 0
}

// IsTransferAccounts 判断当前的消息是不是微信转账
func (m *Message) IsTransferAccounts() bool {
	return m.IsMedia() && m.FileName == "微信转账"
}

// IsSendRedPacket 否发出红包判断当前是
func (m *Message) IsSendRedPacket() bool {
	return m.IsSystem() && m.Content == "发出红包，请在手机上查看"
}

// IsReceiveRedPacket 判断当前是否收到红包
func (m *Message) IsReceiveRedPacket() bool {
	return m.IsSystem() && m.Content == "收到红包，请在手机上查看"
}

// IsRenameGroup 判断当前是否是群组重命名
func (m *Message) IsRenameGroup() bool {
	return m.IsSystem() && strings.Contains(m.Content, "修改群名为")
}

func (m *Message) IsSysNotice() bool {
	return m.MsgType == 9999
}

// StatusNotify 判断是否为操作通知消息
func (m *Message) StatusNotify() bool {
	return m.MsgType == 51
}

// HasFile 判断消息是否为文件类型的消息
func (m *Message) HasFile() bool {
	return m.IsPicture() || m.IsVoice() || m.IsVideo() || (m.IsMedia() && m.AppMsgType == AppMsgTypeAttach) || m.IsEmoticon()
}

// IsSendBySelf 判断消息是否由自己发送
func (m *Message) IsSendBySelf() bool {
	return m.FromUserName == robot.GetBot().GetBotWxId()
}

// IsSendByFriend 判断消息是否由好友发送
func (m *Message) IsSendByFriend() bool {
	return !m.IsSendByGroup() && strings.HasPrefix(m.FromUserName, "@") && !m.IsSendBySelf()
}

// IsSendByGroup 判断消息是否由群组发送
func (m *Message) IsSendByGroup() bool {
	return strings.HasPrefix(m.FromUserName, "@@") || (m.IsSendBySelf() && strings.HasPrefix(m.ToUserName, "@@"))
}

// IsSelfSendToGroup 判断消息是否由自己发送到群组
func (m *Message) IsSelfSendToGroup() bool {
	return m.IsSendBySelf() && strings.HasPrefix(m.ToUserName, "@@")
}
