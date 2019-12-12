package plugin

import (
	"github.com/hanjingo/container"
	"github.com/hanjingo/gate/com"
)

type agentInfoV1 struct {
	id         interface{}     //id
	cache      *container.Pool //缓存
	streamSize int             //可用流量
}
type StreamV1 struct {
	name   string
	agents map[interface{}]*agentInfoV1
}

func NewStreamV1() *StreamV1 {
	return &StreamV1{
		name:   "StreamV1",
		agents: make(map[interface{}]*agentInfoV1),
	}
}

func (s *StreamV1) Name() string {
	return s.name
}

func (s *StreamV1) OnNewAgent(agent com.AgentI) error {
	return nil
}

func (s *StreamV1) OnAgentClose(agent com.AgentI) error {
	return nil
}

func (s *StreamV1) OnMsg(msg *com.Msg) (*com.Msg, error) {
	return msg, nil
}
