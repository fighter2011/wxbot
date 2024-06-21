package ai

import (
	"github.com/sashabaranov/go-openai"
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"sync"
	"time"
)

////go:embed data
//var chatGptData embed.FS

var (
	db          sqlite.DB // 数据库
	chatRoomCtx sync.Map  // 聊天室消息上下文
)

// ChatRoom chatRoomCtx -> ChatRoom => 维系每个人的上下文
type ChatRoom struct {
	chatId   string                         // 聊天室ID, 格式为: 聊天室ID_发送人ID
	chatTime time.Time                      // 聊天时间
	role     string                         // 角色
	content  []openai.ChatCompletionMessage // 聊天上下文内容
}

func init() {
	engine := control.Register("ai", &control.Options{
		Alias:      "AI",
		Help:       "",
		DataFolder: "ai",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/ai.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
}
