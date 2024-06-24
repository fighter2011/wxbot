package robot

import (
	"fmt"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yqchilde/wxbot/engine/pkg/cryptor"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/net"
	"github.com/yqchilde/wxbot/engine/pkg/static"
	"github.com/yqchilde/wxbot/web"
)

var gptdb sqlite.DB

func init() {
	if err := sqlite.Open("data/plugins/chatgpt/chatgpt.db", &gptdb); err != nil {
		log.Fatalf("open sqlite gptdb failed: %v", err)
	}
}

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

	// 设置默认模型
	r.POST("/ai/default-model/upsert/batch", func(c *gin.Context) {
		result := gptdb.Orm.Table("defaultModel").Where("1=1").Delete(&defaultModelDTO{})
		if result.Error != nil {
			c.JSON(http.StatusOK, "清空数据失败")
			return
		}
		var defaultModel []defaultModelDTO
		if err := c.ShouldBindJSON(&defaultModel); err != nil {
			c.JSON(http.StatusOK, "json数组异常")
			return
		}
		result = gptdb.Orm.Table("defaultModel").Create(&defaultModel)
		if result.Error != nil {
			c.JSON(http.StatusOK, "保存数据失败")
			return
		}
		c.JSON(http.StatusOK, "更新成功")
	})
	// 设置模型
	r.POST("/ai/model/upsert/batch", func(c *gin.Context) {
		result := gptdb.Orm.Table("gptmodel").Where("1=1").Delete(&gptModelDTO{})
		if result.Error != nil {
			c.JSON(http.StatusOK, "清空数据失败")
			return
		}
		var gptModel []gptModelDTO
		if err := c.ShouldBindJSON(&gptModel); err != nil {
			c.JSON(http.StatusOK, "json数组异常")
			return
		}
		result = gptdb.Orm.Table("gptmodel").Create(&gptModel)
		if result.Error != nil {
			c.JSON(http.StatusOK, "保存数据失败")
			return
		}
		c.JSON(http.StatusOK, "更新成功")
	})
	// 获取模型列表
	r.GET("/ai/model/list", func(c *gin.Context) {
		var gptModel []gptModelDTO
		if err := gptdb.Orm.Table("gptmodel").Find(&gptModel).Error; err != nil {
			c.JSON(http.StatusOK, "获取数据失败")
			return
		}
		c.JSON(http.StatusOK, gptModel)
	})

	// 获取默认模型列表
	r.GET("/ai/default-model/list", func(c *gin.Context) {
		var gptModel []defaultModelDTO
		if err := gptdb.Orm.Table("defaultModel").Find(&gptModel).Error; err != nil {
			c.JSON(http.StatusOK, "获取数据失败")
			return
		}
		c.JSON(http.StatusOK, gptModel)
	})
	// 更新用户模型列表
	r.POST("/ai/user-model/upsert/batch", func(c *gin.Context) {
		result := gptdb.Orm.Table("userChatModelMapping").Where("1=1").Delete(&userChatModelMappingDTO{})
		if result.Error != nil {
			c.JSON(http.StatusOK, "清空数据失败")
			return
		}
		var userModel []userChatModelMappingDTO
		if err := c.ShouldBindJSON(&userModel); err != nil {
			c.JSON(http.StatusOK, "json数组异常")
			return
		}
		result = gptdb.Orm.Table("userChatModelMapping").Create(&userModel)
		if result.Error != nil {
			c.JSON(http.StatusOK, "保存数据失败")
			return
		}
		c.JSON(http.StatusOK, "更新成功")
	})

	// 获取用户模型列表
	r.GET("/ai/user-model/list", func(c *gin.Context) {
		var userModel []userChatModelMappingDTO
		if err := gptdb.Orm.Table("userChatModelMapping").Find(&userModel).Error; err != nil {
			c.JSON(http.StatusOK, "获取数据失败")
			return
		}
		c.JSON(http.StatusOK, userModel)
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

// gptModelDTO 表名:gptmodel，存放gpt模型相关配置参数 name唯一
type gptModelDTO struct {
	Model            string  `gorm:"column:model" json:"model"`
	Name             string  `gorm:"column:name" json:"name"`
	MaxTokens        int     `gorm:"column:max_tokens" json:"maxTokens"`
	Temperature      float64 `gorm:"column:temperature" json:"temperature"`
	TopP             float64 `gorm:"column:top_p" json:"topP"`
	PresencePenalty  float64 `gorm:"column:presence_penalty" json:"presencePenalty"`
	FrequencyPenalty float64 `gorm:"column:frequency_penalty" json:"frequencyPenalty"`
	ImageSize        string  `gorm:"column:image_size" json:"imageSize"`
	Type             string  `gorm:"column:type" json:"type"`
}

// userChatModelMappingDTO 表名：userChatModelMapping 存放聊天人与模型映射关系
type userChatModelMappingDTO struct {
	ModelName string `gorm:"column:model_name"`
	Uid       string `gorm:"column:uid"`
	Type      string `gorm:"column:type"`
}

// defaultModelDTO 表名:defaultModel 存放默认模型关系
type defaultModelDTO struct {
	ModelName string `gorm:"model_name" json:"modelName"`
	Type      string `gorm:"type" json:"type"` // TEXT/IMAGE/TTS
}
