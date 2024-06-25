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
	"net/http"
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

type model struct {
	name        string // 模型
	displayName string // 展示名称
	supportType string // 支持类型 TEXT/VISION/TTS/IMAGE
}

func getAiClient() (*openai.Client, error) {
	var option AiOption
	if flag, res := redis.Get(AI_PROXY_KEY); flag {
		if len(res) == 0 {
			return nil, errors.New("[AI] AI配置失败")
		}
		if err := json.Unmarshal([]byte(res), &option); err != nil {
			return nil, errors.New("[AI] 解析AI配置失败")
		}
		if len(option.ApiKey) == 0 || len(option.Url) == 0 {
			return nil, errors.New("[AI] 没有有效的AI配置")
		}
	}
	if len(option.Url) == 0 {
		return nil, errors.New("[AI] 未配置proxy")
	}
	config := openai.DefaultConfig(option.ApiKey)
	config.HTTPClient = &http.Client{
		Timeout: time.Minute * 5,
	}
	config.BaseURL = option.Url
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
func getAiModel(uid, supportType string) (*AIModel, error) {
	// 从缓存映射中获取模型
	redisKey := fmt.Sprintf(AI_USER_MODEL_PREFIX_KEY, uid, supportType)
	var aiModel AIModel
	if flag, res := redis.Get(redisKey); flag {
		if err := json.Unmarshal([]byte(res), &aiModel); err != nil {
			return nil, err
		}
	}
	if len(aiModel.Model) > 0 {
		// 存在映射关系 直接返回
		return &aiModel, nil
	}
	// 没有映射从所有数据库里面获取默认模型
	var aiModels []AIModel
	if flag, res := redis.Get(AI_MODEL_KEY); flag {
		if err := json.Unmarshal([]byte(res), &aiModels); err != nil {
			return nil, err
		}
	}
	if len(aiModels) == 0 {
		return nil, errors.New("[AI] 获取AI模型数据失败")
	}

	for _, m := range aiModels {
		if m.IsDefault && m.SupportType == supportType {
			aiModel = m
		}
	}
	if len(aiModel.Model) == 0 {
		return nil, errors.New("[AI] 获取默认模型失败")
	}
	if data, err := json.Marshal(aiModel); err == nil {
		redis.Set(redisKey, string(data))
	} else {
		log.Errorf("[AI] json序列化模型关系异常, %v", aiModel)
	}
	return &aiModel, nil
}

// AskAIWithImage 向ChatGPT请求回复图片
func AskAIWithImage(ctx *robot.Ctx, prompt string, delay ...time.Duration) (b64 string, err error) {
	// 获取客户端
	if aiClient == nil {
		aiClient, err = getAiClient()
		if err != nil {
			return "", err
		}
	}
	aiModel, err := getAiModel(ctx.Uid(), "IMAGE")
	if err != nil {
		return "", err
	}

	// 延迟请求
	if len(delay) > 0 {
		time.Sleep(delay[0])
	}

	resp, err := aiClient.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         prompt,
		Size:           aiModel.ImageSize,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
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
	if len(resp.Data) == 0 {
		return "", ErrServiceUnavailable
	}
	return resp.Data[0].B64JSON, nil
}
