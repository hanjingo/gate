package plugin

import (
	"time"

	"github.com/hanjingo/gate/com"
)

type agentInfoV1 struct {
	id        interface{}     //id
	startTm   time.Time       //开始时间
	subTopics map[string]bool //订阅的主题
	apis      map[uint32]bool //可提供的api服务
}
type ControllerV1 struct {
	name   string
	agents map[interface{}]*agentInfoV1
}

func NewControllerV1() *ControllerV1 {
	return &ControllerV1{
		name:   "ControllerV1",
		agents: make(map[interface{}]*agentInfoV1),
	}
}

func (c *ControllerV1) Name() string {
	return c.name
}

func (c *ControllerV1) OnNewAgent(agent com.AgentI) error {
	return nil
}

func (c *ControllerV1) OnAgentClose(agent com.AgentI) error {
	return nil
}

func (c *ControllerV1) OnMsg(msg *com.Msg) (*com.Msg, error) {
	return msg, nil
}
