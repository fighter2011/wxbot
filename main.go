package main

import (
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/net"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/framework/dean"
	"github.com/yqchilde/wxbot/framework/uos"
	"github.com/yqchilde/wxbot/framework/vlw"
	"time"

	// 导入插件, 变更插件请查看README
	_ "github.com/yqchilde/wxbot/engine/pkg/redis"
	_ "github.com/yqchilde/wxbot/engine/plugins"
)

func main() {

	f := robot.IFramework(nil)
	c := robot.GlobalConfig
	switch c.Framework.Name {
	case "Dean":
		f = robot.IFramework(dean.New(c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping := net.PingConn(ipPort, time.Second*10); !ping {
				c.SetConnHookStatus(false)
				log.Warn("[main] 无法连接Dean框架，网络无法Ping通，请检查网络")
			}
		}
	case "VLW", "vlw":
		f = robot.IFramework(vlw.New(c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
			if ping := net.PingConn(ipPort, time.Second*10); !ping {
				c.SetConnHookStatus(false)
				log.Warn("[main] 无法连接到VLW框架，网络无法Ping通，请检查网络")
			}
		}
	case "UOS":
		f = robot.IFramework(uos.New(c.BotWxId, c.Framework.ApiUrl, c.Framework.ApiToken))
		//if ipPort, err := net.CheckoutIpPort(c.Framework.ApiUrl); err == nil {
		//	if ping := net.PingConn(ipPort, time.Second*10); !ping {
		//		c.SetConnHookStatus(false)
		//		log.Warn("[main] 无法连接到UOS框架，网络无法Ping通，请检查网络")
		//	}
		//}
	default:
		log.Fatalf("[main] 请在配置文件中指定机器人框架后再启动")
	}

	robot.Run(c, f)
}
