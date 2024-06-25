package ai

import (
	"encoding/json"
	"fmt"
	"github.com/yqchilde/wxbot/engine/pkg/redis"
	"github.com/yqchilde/wxbot/engine/robot"
	"regexp"
)

// 设置模型相关指令
func setModelCommand(ctx *robot.Ctx, msg string, command string) {
	switch command {
	case "模型列表":
		replyMsg := "模型列表:\n"
		var aiModels []AIModel
		if f, res := redis.Get(AI_MODEL_KEY); f {
			if err := json.Unmarshal([]byte(res), &aiModels); err != nil {
				ctx.ReplyTextAndAt("获取模型列表失败")
				return
			}
			if len(aiModels) == 0 {
				ctx.ReplyTextAndAt("未配置模型列表")
				return
			}
		} else {
			ctx.ReplyTextAndAt("解析模型列表失败")
			return
		}
		for _, m := range aiModels {
			replyMsg += fmt.Sprintf("%s\n", m.Name)
		}
		ctx.ReplyTextAndAt(replyMsg)
	case "切换模型":
		modelName := regexp.MustCompile(`^切换模型\s+(.+)$`).FindStringSubmatch(msg)[1]
		if len(modelName) == 0 {
			ctx.ReplyTextAndAt("模型不存在")
			return
		}
		var aiModels []AIModel
		if f, res := redis.Get(AI_MODEL_KEY); f {
			if err := json.Unmarshal([]byte(res), &aiModels); err != nil {
				ctx.ReplyTextAndAt("切换模型失败")
				return
			}
			if len(aiModels) == 0 {
				ctx.ReplyTextAndAt("切换模型失败")
				return
			}
		} else {
			ctx.ReplyTextAndAt("切换模型失败")
			return
		}
		redisKey := fmt.Sprintf(AI_USER_MODEL_PREFIX_KEY, ctx.Uid(), "TEXT")
		for _, m := range aiModels {
			if modelName == m.Model {
				data, err := json.Marshal(m)
				if err != nil {
					ctx.ReplyTextAndAt("切换模型失败")
					return
				}
				if redis.Set(redisKey, string(data)) {
					ctx.ReplyTextAndAt("切换模型成功")
				}
			}
		}
	case "当前模型":
		redisKey := fmt.Sprintf(AI_USER_MODEL_PREFIX_KEY, ctx.Uid(), "TEXT")
		var aiModel AIModel
		if f, res := redis.Get(redisKey); f {
			if err := json.Unmarshal([]byte(res), &aiModel); err != nil {
				ctx.ReplyTextAndAt("获取模型失败")
				return
			}
			if len(aiModel.Model) == 0 {
				ctx.ReplyTextAndAt("模型为空")
				return
			}
			ctx.ReplyTextAndAt("当前模型为:\n" + aiModel.Name)
		} else {
			ctx.ReplyTextAndAt("获取模型失败")
			return
		}
	}

}
