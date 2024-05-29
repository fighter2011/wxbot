package uos

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
	"net/http"
)

const (
	MsgTypeText           int = 1     // 文本消息
	MsgTypeImage          int = 3     // 图片消息
	MsgTypeVoice          int = 34    // 语音消息
	MsgTypeVerify         int = 37    // 认证消息
	MsgTypePossibleFriend int = 40    // 好友推荐消息
	MsgTypeShareCard      int = 42    // 名片消息
	MsgTypeVideo          int = 43    // 视频消息
	MsgTypeEmoticon       int = 47    // 表情消息
	MsgTypeLocation       int = 48    // 地理位置消息
	MsgTypeApp            int = 49    // APP消息
	MsgTypeVoip           int = 50    // VOIP消息
	MsgTypeVoipNotify     int = 52    // VOIP结束消息
	MsgTypeVoipInvite     int = 53    // VOIP邀请
	MsgTypeMicroVideo     int = 62    // 小视频消息
	MsgTypeSys            int = 10000 // 系统消息
	MsgTypeRecalled       int = 10002 // 消息撤回
)

type Framework struct {
	BotWxId  string // 机器人微信ID
	ApiUrl   string // http api地址
	ApiToken string // http api鉴权token
}

func New(botWxId string) *Framework {
	return &Framework{}
}

func (f *Framework) Callback(ctx *gin.Context, handler func(*robot.Event, robot.IFramework)) {
	recv, err := ctx.GetRawData()
	if err != nil {
		log.Errorf("[UOS] 接收回调错误, error: %v", err)
		return
	}
	handler(buildEvent(recv), f)
	ctx.JSON(http.StatusOK, gin.H{"code": 0})
}

func buildEvent(resp []byte) *robot.Event {
	var event robot.Event
	var msg Message
	err := json.Unmarshal(resp, &msg)
	if err != nil {
		log.Errorf("解析消息失败 %v 异常是: %v", string(resp), err)
		//todo notify
		return nil
	}

	switch msg.MsgType {
	//case eventAccountChange:
	// todo
	case MsgTypeText:
		log.Printf("收到文本消息, %v", msg)
		//	switch gjson.Get(resp, "data.fromType").Int() {
		//	case 1: // 私聊
		//		switch gjson.Get(resp, "data.msgType").Int() {
		//		case 49: // 私聊发app应用消息
		//			event = robot.Event{
		//				Type:         robot.EventPrivateChat,
		//				FromUniqueID: gjson.Get(resp, "data.fromWxid").String(),
		//				FromWxId:     gjson.Get(resp, "data.fromWxid").String(),
		//				FromName:     "",
		//				Message: &robot.Message{
		//					Type:    gjson.Get(resp, "data.msgType").Int(),
		//					Content: gjson.Get(resp, "data.msg").String(),
		//				},
		//			}
		//
		//			//var refer ReferenceXml
		//			//if err := xml.Unmarshal([]byte(gjson.Get(resp, "data.msg").String()), &refer); err == nil {
		//			//	if refer.Appmsg.Refermsg != nil { // 引用消息
		//			//		event.Message.Type = robot.MsgTypeText // 方便匹配
		//			//		event.Message.Content = refer.Appmsg.Title
		//			//		event.ReferenceMessage = &robot.ReferenceMessage{
		//			//			FromUser:    refer.Appmsg.Refermsg.Fromusr,
		//			//			ChatUser:    refer.Appmsg.Refermsg.Chatusr,
		//			//			DisplayName: refer.Appmsg.Refermsg.Displayname,
		//			//			Content:     refer.Appmsg.Refermsg.Content,
		//			//		}
		//			//	}
		//			//}
		//		default:
		//			event = robot.Event{
		//				Type:         robot.EventPrivateChat,
		//				FromUniqueID: gjson.Get(resp, "data.fromWxid").String(),
		//				FromWxId:     gjson.Get(resp, "data.fromWxid").String(),
		//				FromName:     "",
		//				IsAtMe:       true,
		//				Message: &robot.Message{
		//					Type:    gjson.Get(resp, "data.msgType").Int(),
		//					Content: gjson.Get(resp, "data.msg").String(),
		//				},
		//			}
		//			for _, data := range robot.GetBot().Friends() {
		//				if data.WxId == event.FromWxId {
		//					event.FromName = data.Nick
		//					event.FromUniqueName = data.Nick
		//					break
		//				}
		//			}
		//		}
		//	case 2: // 群聊
		//		switch gjson.Get(resp, "data.msgType").Int() {
		//		case 10000:
		//			event = robot.Event{
		//				Type: robot.EventSystem,
		//				Message: &robot.Message{
		//					Content: gjson.Get(resp, "data.msg").String(),
		//				},
		//			}
		//		case 49: // 群聊发app应用消息
		//			event = robot.Event{
		//				Type:          robot.EventGroupChat,
		//				FromUniqueID:  gjson.Get(resp, "data.fromWxid").String(),
		//				FromGroup:     gjson.Get(resp, "data.fromWxid").String(),
		//				FromGroupName: "",
		//				FromWxId:      gjson.Get(resp, "data.finalFromWxid").String(),
		//				FromName:      "",
		//				Message: &robot.Message{
		//					Type:    gjson.Get(resp, "data.msgType").Int(),
		//					Content: gjson.Get(resp, "data.msg").String(),
		//				},
		//			}
		//
		//			//var refer ReferenceXml
		//			//if err := xml.Unmarshal([]byte(gjson.Get(resp, "data.msg").String()), &refer); err == nil {
		//			//	if refer.Appmsg.Refermsg != nil { // 引用消息
		//			//		event.Message.Type = robot.MsgTypeText // 方便匹配
		//			//		event.Message.Content = refer.Appmsg.Title
		//			//		event.ReferenceMessage = &robot.ReferenceMessage{
		//			//			FromUser:    refer.Appmsg.Refermsg.Fromusr,
		//			//			ChatUser:    refer.Appmsg.Refermsg.Chatusr,
		//			//			DisplayName: refer.Appmsg.Refermsg.Displayname,
		//			//			Content:     refer.Appmsg.Refermsg.Content,
		//			//		}
		//			//	}
		//			//}
		//		default:
		//			event = robot.Event{
		//				Type:          robot.EventGroupChat,
		//				FromUniqueID:  gjson.Get(resp, "data.fromWxid").String(),
		//				FromGroup:     gjson.Get(resp, "data.fromWxid").String(),
		//				FromGroupName: "",
		//				FromWxId:      gjson.Get(resp, "data.finalFromWxid").String(),
		//				FromName:      "",
		//				Message: &robot.Message{
		//					Type:    gjson.Get(resp, "data.msgType").Int(),
		//					Content: gjson.Get(resp, "data.msg").String(),
		//				},
		//			}
		//			if gjson.Get(resp, fmt.Sprintf("data.atWxidList.#(==%s)", gjson.Get(resp, "wxid").String())).Exists() {
		//				if !strings.Contains(event.Message.Content, "@所有人") {
		//					event.IsAtMe = true
		//				}
		//			}
		//			for _, data := range robot.GetBot().Groups() {
		//				if data.WxId == event.FromGroup {
		//					event.FromGroupName = data.Nick
		//					event.FromUniqueName = data.Nick
		//					break
		//				}
		//			}
		//		}
		//	case 3: // 公众号
		//		event = robot.Event{
		//			Type:         robot.EventMPChat,
		//			FromUniqueID: gjson.Get(resp, "data.fromWxid").String(),
		//			FromWxId:     gjson.Get(resp, "data.fromWxid").String(),
		//			FromName:     "",
		//			MPMessage: &robot.Message{
		//				Type:    gjson.Get(resp, "data.msgType").Int(),
		//				Content: gjson.Get(resp, "data.msg").String(),
		//			},
		//		}
		//		for _, data := range robot.GetBot().MPs() {
		//			if data.WxId == event.FromWxId {
		//				event.FromName = data.Nick
		//				event.FromUniqueName = data.Nick
		//				break
		//			}
		//		}
		//	}
		//
		//	// 自身发言
		////case eventSelfMessage:
		////	event = robot.Event{
		////		Type:         robot.EventSelfMessage,
		////		FromUniqueID: gjson.Get(resp, "data.data.fromWxid").String(),
		////		FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
		////		Message: &robot.Message{
		////			Type:    gjson.Get(resp, "data.data.msgType").Int(),
		////			Content: gjson.Get(resp, "data.data.msg").String(),
		////		},
		////	}
		//case eventTransfer:
		//	event = robot.Event{
		//		Type: robot.EventTransfer,
		//		TransferMessage: &robot.TransferMessage{
		//			FromWxId:     gjson.Get(resp, "data.fromWxid").String(),
		//			MsgSource:    gjson.Get(resp, "data.msgSource").Int(),
		//			TransferType: gjson.Get(resp, "data.transType").Int(),
		//			Money:        gjson.Get(resp, "data.money").String(),
		//			Memo:         gjson.Get(resp, "data.memo").String(),
		//			TransferId:   gjson.Get(resp, "data.transferid").String(),
		//			TransferTime: gjson.Get(resp, "data.invalidtime").String(),
		//		},
		//	}
		//case eventMessageWithdraw:
		//	fromType := gjson.Get(resp, "data.fromType").Int()
		//	if fromType == 1 {
		//		event = robot.Event{
		//			Type: robot.EventMessageWithdraw,
		//			WithdrawMessage: &robot.WithdrawMessage{
		//				FromType:  fromType,
		//				FromWxId:  gjson.Get(resp, "data.fromWxid").String(),
		//				MsgSource: gjson.Get(resp, "data.msgSource").Int(),
		//				Msg:       gjson.Get(resp, "data.msg").String(),
		//			},
		//		}
		//	} else if fromType == 2 {
		//		event = robot.Event{
		//			Type: robot.EventMessageWithdraw,
		//			WithdrawMessage: &robot.WithdrawMessage{
		//				FromType:  fromType,
		//				FromGroup: gjson.Get(resp, "data.fromWxid").String(),
		//				FromWxId:  gjson.Get(resp, "data.finalFromWxid").String(),
		//				MsgSource: gjson.Get(resp, "data.msgSource").Int(),
		//				Msg:       gjson.Get(resp, "data.msg").String(),
		//			},
		//		}
		//	}
		//case eventFriendVerify:
		//	event = robot.Event{
		//		Type: robot.EventFriendVerify,
		//		FriendVerifyMessage: &robot.FriendVerifyMessage{
		//			WxId:      gjson.Get(resp, "data.wxid").String(),
		//			Nick:      gjson.Get(resp, "data.nick").String(),
		//			V3:        gjson.Get(resp, "data.v3").String(),
		//			V4:        gjson.Get(resp, "data.v4").String(),
		//			AvatarUrl: gjson.Get(resp, "data.avatarMinUrl").String(),
		//			Content:   gjson.Get(resp, "data.content").String(),
		//			Scene:     gjson.Get(resp, "data.scene").String(),
		//		},
		//	}
	}

	//event.RobotWxId = gjson.Get(resp, "wxid").String()
	//event.RawMessage = resp
	return &event
}
