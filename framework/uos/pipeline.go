package uos

import (
	"errors"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

func newPipeline() *Pipeline {
	return &Pipeline{}
}

type Pipeline struct {
	matchNodes matchNodes
}

// MatchFunc 消息匹配函数,返回为true则表示匹配
type MatchFunc func(*Message) bool

// Processor 消息处理函数
type Processor func(*Message) *robot.Event

type matchNode struct {
	matchFunc MatchFunc
	processor Processor
}

type matchNodes []*matchNode

func (p *Pipeline) RegisterProcessor(matchFun MatchFunc, processors Processor) {
	p.matchNodes = append(p.matchNodes, &matchNode{matchFunc: matchFun, processor: processors})
}

// Process
// 获取匹配方法进行执行
func (p *Pipeline) doProcess(msg *Message) (*robot.Event, error) {
	var processors []Processor
	for _, node := range p.matchNodes {
		if node.matchFunc(msg) {
			processors = append(processors, node.processor)
		}
	}
	if len(processors) == 0 {
		return nil, errors.New("not Found processor to process")
	}
	if len(processors) > 1 {
		log.Warnf("processor more than 1, %v", msg)
	}

	return processors[0](msg), nil
}
