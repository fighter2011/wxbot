package robot

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/yqchilde/wxbot/engine/pkg/redis"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yqchilde/wxbot/engine/pkg/cryptor"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/net"
	"github.com/yqchilde/wxbot/engine/pkg/static"
	"github.com/yqchilde/wxbot/web"
)

const (
	AI_PROXY_KEY             string = "ai:one-api:proxy"
	AI_MODEL_KEY             string = "ai:model:list"
	AI_USER_MODEL_PREFIX_KEY string = "ai:model:user:%s:%s" // 用户model映射key accountId + supportType
	AI_ROLE_KEY              string = "ai:role:list"
)

// 跨域 middleware
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func runServer(c *Config) {
	r := gin.New()
	r.Use(cors())
	r.Use(static.Serve("/", static.EmbedFolder(web.Web, "dist")))

	// 消息回调
	r.POST("/wxbot/callback", func(c *gin.Context) {
		bot.framework.Callback(c, eventBuffer.ProcessEvent)
	})

	// 静态文件服务
	r.GET("/wxbot/static", func(c *gin.Context) {
		if c.Query("file") == "" {
			c.String(http.StatusInternalServerError, "Warning: 非法访问")
			return
		}
		filename, err := cryptor.DecryptFilename(fileSecret, c.Query("file"))
		if err != nil {
			log.Errorf("[http] 静态文件解密失败: %s", err.Error())
			c.String(http.StatusInternalServerError, "Warning: 非法访问")
			return
		}
		if !strings.HasPrefix(filename, "data/plugins") && !strings.HasPrefix(filename, "./data/plugins") &&
			!strings.HasPrefix(filename, "data\\plugins") && !strings.HasPrefix(filename, ".\\data\\plugins") {
			log.Errorf("[http] 非法访问静态文件: %s", filename)
			c.String(http.StatusInternalServerError, "Warning: 非法访问")
			return
		}
		c.File(filename)
	})

	// 菜单接口
	r.GET("/wxbot/menu", func(c *gin.Context) {
		wxId := c.Query("wxid")
		if wxId == "" || wxId == "undefined" {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "wxid不能为空",
			})
			return
		}

		menus := ControlApi.GetMenus(wxId)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": menus,
		})
	})

	r.GET("/wxbot/ai/model/list", func(c *gin.Context) {
		var aiModels []AIModel
		if f, res := redis.Get(AI_MODEL_KEY); f {
			if err := json.Unmarshal([]byte(res), &aiModels); err != nil {
				c.JSON(http.StatusOK, gin.H{})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": aiModels})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
		return
	})

	r.POST("/wxbot/ai/model/upsert/batch", func(c *gin.Context) {
		var aiModels []AIModel
		if err := c.ShouldBindJSON(&aiModels); err != nil {
			c.JSON(http.StatusOK, "json数组异常")
			return
		}
		data, err := json.Marshal(&aiModels)
		if err != nil {
			c.JSON(http.StatusOK, "json数组异常")
			return
		}
		if flag := redis.Set(AI_MODEL_KEY, string(data)); flag {
			c.JSON(http.StatusOK, "设置模型成功")
			return
		}
	})

	// no route
	r.NoRoute(func(c *gin.Context) {
		c.FileFromFS("/", static.EmbedFolder(web.Web, "dist"))
	})

	if ip, err := net.GetIPWithLocal(); err != nil {
		log.Printf("[robot] WxBot回调地址: http://%s:%d/wxbot/callback", "127.0.0.1", c.ServerPort)
	} else {
		log.Printf("[robot] WxBot回调地址: http://%s:%d/wxbot/callback", ip, c.ServerPort)
	}
	if err := r.Run(fmt.Sprintf(":%d", c.ServerPort)); err != nil {
		log.Fatalf("[robot] WxBot回调服务启动失败, error: %v", err)
	}
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
