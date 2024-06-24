package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/redis"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/plugins/chatgpt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	aiClient *openai.Client
)

var (
	ErrNoKey              = fmt.Errorf("请先私聊机器人配置apiKey\n指令：set chatgpt apikey __(多个key用;符号隔开)\napiKey获取请到https://beta.openai.com获取")
	ErrMaxTokens          = errors.New("OpenAi免费上下文长度限制为4097个词组，您的上下文长度已超出限制")
	ErrExceededQuota      = errors.New("OpenAi配额已用完，请联系管理员")
	ErrIncorrectKey       = errors.New("OpenAi ApiKey错误，请联系管理员")
	ErrServiceUnavailable = errors.New("ChatGPT服务异常，请稍后再试")
)

type aiOption struct {
	apiKey string  // apiKey
	url    string  // 请求地址
	models []model // 支持模型列表
}

type model struct {
	name        string // 模型
	displayName string // 展示名称
	supportType string // 支持类型 TEXT/VISION/TTS/IMAGE
}

func getAiClient() (*openai.Client, error) {
	var option aiOption
	if flag, res := redis.Get("ai:one-api:key"); flag {
		if len(res) == 0 {
			return nil, errors.New("[AI] AI配置失败")
		}
		if err := json.Unmarshal([]byte(res), &option); err != nil {
			return nil, errors.New("[AI] 解析AI配置失败")
		}
		if len(option.apiKey) == 0 || len(option.url) == 0 || len(option.models) == 0 {
			return nil, errors.New("[AI] 没有有效的AI配置")
		}
	}
	proxyUrl, err := url.Parse(option.url)
	if err != nil {
		log.Errorf("[AI] 解析http_proxy失败, error:%s", err.Error())
		return nil, errors.New("解析http_proxy失败")
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config := openai.DefaultConfig(option.apiKey)
	config.HTTPClient = &http.Client{
		Transport: transport,
		Timeout:   time.Minute * 5,
	}
	return openai.NewClientWithConfig(config), nil
}

// AskAI 向AI进行请求回复
func AskAI(ctx *robot.Ctx, messages []openai.ChatCompletionMessage, delay ...time.Duration) (answer string, err error) {
	// 获取客户端
	if aiClient == nil {
		aiClient, err = getAiClient()
		if err != nil {
			return "", err
		}
	}
	aiModel, err := getAiModel(ctx.Uid(), "TEXT")
	if err != nil {
		return "", err
	}
	if len(delay) > 0 {
		time.Sleep(delay[0])
	}
	// 处理用户role
	var role string
	if val, ok := roomCtx.Load(ctx.Uid()); ok {
		role = val.(Room).role
	}
	if len(role) == 0 {
		role = "默认"
	}
	var chatMessages []openai.ChatCompletionMessage
	if strings.Contains(SystemRole.MustGet(role).(string), "%s") {
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    "system",
			Content: fmt.Sprintf(SystemRole.MustGet(role).(string), robot.GetBot().GetConfig().BotNickname),
		})
	} else {
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    "system",
			Content: SystemRole.MustGet(role).(string),
		})
	}
	chatMessages = append(chatMessages, messages...)
	resp, err := aiClient.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    aiModel.Model,
		Messages: chatMessages,
	})
	// 处理响应回来的错误
	if err != nil {
		if strings.Contains(err.Error(), "Please reduce your prompt") || strings.Contains(err.Error(), "Please reduce the length of the messages") {
			return "", ErrMaxTokens
		}
		if strings.Contains(err.Error(), "Incorrect API key") {
			return "", ErrIncorrectKey
		}
		if strings.Contains(err.Error(), "invalid character") {
			return "", ErrServiceUnavailable
		}
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", ErrServiceUnavailable
	}
	//todo 超过2000字的回复 生成图片进行回复
	return resp.Choices[0].Message.Content, nil
}

// 获取使用模型
func getAiModel(uid, supportType string) (*chatgpt.GptModel, error) {

	return nil, nil
}
