package ai

import (
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/utils"
	"github.com/yqchilde/wxbot/engine/robot"
	"path/filepath"
	"time"
)

// 设置图片相关指令
func setImageCommand(ctx *robot.Ctx, msg string, command string) {
	switch command {
	case "作画":
		b64, err := AskAIWithImage(ctx, msg, time.Second)
		if err != nil {
			log.Errorf("AI出错了，Err：%s", err.Error())
			ctx.ReplyTextAndAt("AI出错了，Err：" + err.Error())
			return
		}
		filename := filepath.Join("data/plugins/chatgpt/cache", msg+".png")
		if err = utils.Base64ToImage(b64, filename); err != nil {
			log.Errorf("作画失败，Err: %s", err.Error())
			ctx.ReplyTextAndAt("作画失败，请重试")
			return
		}
		ctx.ReplyImage("local://" + filename)
	}
}

// 设置单次提问指令
func setSingleCommand(ctx *robot.Ctx, msg string, command string) {
	switch command {
	case "提问":
		messages := []openai.ChatCompletionMessage{{Role: "user", Content: msg}}
		answer, err := AskAI(ctx, messages, time.Second)
		if err != nil {
			if errors.Is(err, ErrNoKey) {
				ctx.ReplyTextAndAt(err.Error())
			} else {
				ctx.ReplyTextAndAt("ChatGPT出错了，Err：" + err.Error())
			}
			return
		}
		//answer = replaceSensitiveWords(answer)
		ctx.ReplyTextAndAt(fmt.Sprintf("问：%s \n--------------------\n答：%s", msg, answer))
	}
}
