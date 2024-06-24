package ai

import (
	"errors"
	"github.com/sashabaranov/go-openai"
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
	"sync"
	"time"
)

var (
	roomCtx sync.Map //聊天室上下文
)

type Room struct {
	id       string                         // 聊天室id 格式为：群id/发送者id_发送者id
	chatTime time.Time                      // 聊天时间
	role     string                         //角色
	content  []openai.ChatCompletionMessage // 聊天上下文
}

func init() {
	engine := control.Register("ai", &control.Options{
		Alias: "AI",
		Help: "指令:\n" +
			"* @机器人 [内容] -> 进行AI对话，计入上下文\n" +
			"* @机器人 提问 [问题] -> 单独提问，不计入上下文\n" +
			"* @机器人 作画 [描述] -> 进行AI作画\n" +
			"* @机器人 清空会话 -> 可清空与您的上下文\n" +
			"* @机器人 角色列表 -> 获取可切换的AI角色\n" +
			"* @机器人 当前角色 -> 获取当前用户的AI角色\n" +
			"* @机器人 创建角色 [角色名] [角色描述]\n" +
			"* @机器人 删除角色 [角色名]\n" +
			"* @机器人 切换角色 [角色名]\n\n" +
			"*管理员指令(详细说明请看文档):\n" +
			"* set chatgpt apikey [keys]\n" +
			"* del chatgpt apikey [keys]\n" +
			"* set chatgpt model [key=val]\n" +
			"* reset chatgpt model\n" +
			"* get chatgpt info\n" +
			"* set chatgpt proxy [url]\n" +
			"* del chatgpt proxy\n" +
			"* set chatgpt http_proxy [url]\n" +
			"* del chatgpt http_proxy\n" +
			"* get chatgpt (sensitive|敏感词)\n" +
			"* set chatgpt (sensitive|敏感词) [敏感词]\n" +
			"* reset chatgpt (sensitive|敏感词)\n" +
			"* del chatgpt system (sensitive|敏感词)\n" +
			"* del chatgpt user (sensitive|敏感词)\n" +
			"* del chatgpt all (sensitive|敏感词)",
		DataFolder: "ai",
	})

	engine.OnMessage(robot.OnlyAtMe).SetBlock(true).SetPriority(9999).Handle(func(ctx *robot.Ctx) {
		var (
			now  = time.Now().Local()
			msg  = ctx.MessageString()
			room = Room{
				id:       ctx.Uid(),
				chatTime: now,
				content:  []openai.ChatCompletionMessage{},
			}
		)

		//todo 敏感词过滤
		//todo 预判断 即指令处理 此处可以考虑command/策略模式

		//正式处理开始
		if c, ok := roomCtx.Load(msg); ok {
			// 判断距离上次聊天是否超过10分钟了
			if now.Sub(c.(Room).chatTime) > 10*time.Minute {
				roomCtx.LoadAndDelete(room.id)
				room.content = []openai.ChatCompletionMessage{{Role: "user", Content: msg}}
			} else {
				room.content = append(c.(Room).content, openai.ChatCompletionMessage{Role: "user", Content: msg})
			}
		} else {
			room.content = []openai.ChatCompletionMessage{{Role: "user", Content: msg}}
		}
		answer, err := AskChatGpt(ctx, chatRoom.content, time.Second)
		if err != nil {
			switch {
			case errors.Is(err, ErrNoKey):
				ctx.ReplyTextAndAt(err.Error())
			case errors.Is(err, ErrMaxTokens):
				ctx.ReplyTextAndAt("和你的聊天上下文内容太多啦，我的记忆好像在消退.. 糟糕，我忘记了..，请重新问我吧")
				roomCtx.LoadAndDelete(room.id)
			default:
				ctx.ReplyTextAndAt("AI出错了，Err：" + err.Error())
			}
			return
		}

		room.content = append(room.content, openai.ChatCompletionMessage{Role: "assistant", Content: answer})
		roomCtx.Store(room.id, room)
		answer = replaceSensitiveWords(answer)
		ctx.ReplyTextAndAt(answer)
	})

}
