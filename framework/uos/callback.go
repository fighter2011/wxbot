package uos

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
	"net/http"
)

type Framework struct {
	BotWxId  string // 机器人微信ID
	ApiUrl   string // http api地址
	ApiToken string // http api鉴权token
	pipeline *Pipeline
}

func New(botWxId, apiUrl, apiToken string) *Framework {
	return &Framework{
		BotWxId:  botWxId,
		ApiUrl:   apiUrl,
		ApiToken: apiToken,
		pipeline: newPipeline(),
	}
}
func (f *Framework) Callback(ctx *gin.Context, handler func(*robot.Event, robot.IFramework)) {
	recv, err := ctx.GetRawData()
	if err != nil {
		log.Errorf("[UOS] 接收回调错误, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": fmt.Sprintf("%s:%s", "读取消息失败", err.Error())})
		return
	}
	log.Printf("收到回调消息内容:\n %v", string(recv))
	var event *robot.Event
	var callbackResp CallbackResp
	err = json.Unmarshal(recv, &callbackResp)
	if err != nil {
		log.Errorf("解析消息失败 %v 异常是: %v", string(recv), err)
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": fmt.Sprintf("%s:%s", "解析json失败", err.Error())})
		return
	}
	event, err = f.pipeline.doProcess(&callbackResp.Data)
	if err != nil {
		log.Errorf("解析消息失败 %v 异常是: %v", string(recv), err)
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": fmt.Sprintf("%s:%s", "处理消息失败", err.Error())})
		return
	}
	handler(event, f)
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}

func (f *Framework) Init() {
	// 处理私聊文字/图片消息
	f.pipeline.RegisterProcessor(func(msg *Message) bool {
		return msg.IsSendByFriend() && (msg.IsText() || msg.HasFile())
	}, func(msg *Message) *robot.Event {
		var event *robot.Event
		var content string
		if msg.IsUrl() {
			content = msg.AttachmentUrl
		} else {
			content = msg.Content
		}
		event = &robot.Event{
			IsAtMe:       true,
			Type:         robot.EventPrivateChat,
			FromUniqueID: msg.FromUserName,
			FromWxId:     msg.FromUserName,
			FromName:     "",
			Message: &robot.Message{
				Id:      msg.MsgId,
				Type:    int64(msg.MsgType),
				Content: content,
			},
		}

		for _, data := range robot.GetBot().Friends() {
			if data.WxId == event.FromWxId {
				event.FromName = data.Nick
				break
			}
		}
		return event
	})
	// 处理添加好友消息
	f.pipeline.RegisterProcessor(func(msg *Message) bool {
		return msg.IsFriendAdd()
	}, func(msg *Message) *robot.Event {
		var event *robot.Event
		event = &robot.Event{
			Type:         robot.EventPrivateChat,
			FromUniqueID: msg.FromUserName,
			FromWxId:     msg.FromUserName,
			FromName:     "",
			Message: &robot.Message{
				Type:    int64(msg.MsgType),
				Content: msg.Content,
			},
			FriendVerifyMessage: &robot.FriendVerifyMessage{
				WxId:          msg.RecommendInfo.UserName,
				Nick:          msg.RecommendInfo.NickName,
				Content:       msg.RecommendInfo.Content,
				RecommendInfo: msg.RecommendInfo,
			},
		}
		return event
	})
	// 处理群聊文字
	f.pipeline.RegisterProcessor(func(msg *Message) bool {
		return msg.IsSendByGroup() && (msg.IsText() || msg.HasFile())
	}, func(msg *Message) *robot.Event {
		var event *robot.Event
		var content string
		if msg.IsUrl() {
			content = msg.AttachmentUrl
		} else {
			content = msg.Content
		}
		event = &robot.Event{
			Type:          robot.EventGroupChat,
			FromUniqueID:  msg.FromUserName,
			FromGroup:     msg.FromUserName,
			FromGroupName: "",
			FromName:      "",
			IsAtMe:        msg.IsAt,
			Message: &robot.Message{
				Type:    int64(msg.MsgType),
				Content: content,
			},
		}
		// 设置发送用户信息
		if msg.GroupSenderInfo.WxId != "" {
			event.FromWxId = msg.GroupSenderInfo.WxId
			event.FromName = msg.GroupSenderInfo.NickName
		}

		for _, data := range robot.GetBot().Groups() {
			if data.WxId == event.FromGroup {
				event.FromGroupName = data.Nick
				event.FromUniqueName = data.Nick
				break
			}
		}
		return event
	})
	// 处理转帐消息
	f.pipeline.RegisterProcessor(func(msg *Message) bool {
		return false
	}, func(msg *Message) *robot.Event {
		return nil
	})
}
