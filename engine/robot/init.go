package robot

import (
	"github.com/spf13/viper"
	"github.com/yqchilde/wxbot/engine/pkg/database"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

var configTemplate = `# 机器人WxId，修改为自己的机器人wxId
botWxId: "你的机器人wxId"
# 机器人名字
botNickname: "Q宝"
# 管理员wxId，多个管理员依次添加，用于管理机器人的wxId
superUsers:
  - "你自己的wxId"
# 管理员命令前缀，匹配系统内置管理员指令需要
commandPrefix: "/"
# 唤醒机器人限制，可选:at，表示群聊必须at机器人才可匹配指令，留空时为默认规则(无特殊要求建议留空)
wakeUpRequire: ""

# 本项目运行时会启动一个HTTP服务，包含一个接收事件回调服务，该配置项为HTTP服务端口
# 在所接入的VX框架中请将回调地址改为 http://[本项目运行服务器IP]:[serverPort]/wxbot/callback
serverPort: 9528
# 本项目运行时会启动一个HTTP服务，包含一个静态图片文件服务，用于将本地图片作为网络图片
# 仅当在插件中使用ctx.ReplyImage(local://[图片路径时])才会用到该项
# 本项目和VX框架运行在同一台服务器时，请填写http://[本机IP]:[serverPort]，否则请填写本项目运行服务器IP，也可以使用域名
serverAddress: "http://192.168.31.12:9528"

# 接入框架配置
framework:
  # 框架选择，可选 Dean、VLW
  name: "Dean"
  # wxbot主动请求微信框架的地址，比如：http://[运行Dean的服务器ip]:[Dean的HTTP端口]
  apiUrl: "http://192.168.31.8:9527"
  # VX框架HTTP鉴权Token (Dean目前没有，vlw需要)
  apiToken: ""
`

var version string
var GlobalConfig *Config

func init() {
	// 检查配置文件是否存在
	if _, err := os.Stat("./config/config.yaml"); os.IsNotExist(err) {
		log.Println("未发现配置文件，已为您生成配置文件，请修改后重新运行程序")
		if err := os.WriteFile("config.yaml", []byte(configTemplate), 0644); err != nil {
			log.Fatalf("生成配置文件失败: %v", err)
		}
		os.Exit(0)
	}
	v := viper.New()
	v.SetConfigFile("config/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("[main] 读取配置文件失败: %s", err.Error())
	}
	GlobalConfig = NewConfig()
	if err := v.Unmarshal(GlobalConfig); err != nil {
		log.Fatalf("[main] 解析配置文件失败: %s", err.Error())
	}
	// 初始化redis
	initRedis(GlobalConfig)
	// 打印版本
	log.Printf("当前运行版本: %s", version)
	gin.SetMode(gin.ReleaseMode)
}

func initRedis(c *Config) {
	database.InitRedisConn(&database.RedisOption{
		Addr:     c.Redis.Addr,
		DB:       c.Redis.DB,
		Password: c.Redis.Password,
		PoolSize: c.Redis.PoolSize,
	})
}
