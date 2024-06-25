package ai

import "github.com/sashabaranov/go-openai"

// SystemRoles 表名:roles，存放系统角色
type SystemRoles struct {
	Role string `gorm:"column:role"`
	Desc string `gorm:"column:desc"`
}

type AiOption struct {
	ApiKey string `json:"apiKey,omitempty"` // apiKey
	Url    string `json:"url,omitempty"`    // 请求地址
}

type AIModel struct {
	Model            string            `json:"model"` // model
	Name             string            `json:"name"`  // 展示名称 仅展示
	MaxTokens        int               `json:"maxTokens"`
	Temperature      float32           `json:"temperature"`
	TopP             float32           `json:"topP"`
	PresencePenalty  float32           `json:"presencePenalty"`
	FrequencyPenalty float32           `json:"frequencyPenalty"`
	ImageSize        string            `json:"imageSize"`
	ToolCalls        []openai.ToolCall `json:"toolCalls"` // 冗余 等simple-one-api支持
	ToolCallID       string            `json:"toolCallID"`
	SupportType      string            `json:"supportType"` // 支持类型 TEXT/IMAGE/TTS
	IsDefault        bool              `json:"isDefault"`   // 默认类型
}
