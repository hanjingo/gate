package plugin

import (
	"time"

	"github.com/hanjingo/gate/com"
	pc4 "github.com/hanjingo/protocol/v4"
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
	codec  *pc4.Codec
}

func NewControllerV1() *ControllerV1 {
	return &ControllerV1{
		name:   "ControllerV1",
		agents: make(map[interface{}]*agentInfoV1),
		codec:  pc4.NewCodec(),
	}
}

func (c *ControllerV1) Name() string {
	return c.name
}

func (c *ControllerV1) OnNewAgent(agent com.AgentI) error {
	info := &agentInfoV1{
		id:        agent.GetId(),
		startTm:   time.Now(),
		subTopics: make(map[string]bool),
		apis:      make(map[uint32]bool),
	}
	c.agents[info.id] = info
	return nil
}

func (c *ControllerV1) OnAgentClose(agent com.AgentI) error {
	if info, ok := c.agents[agent.GetId()]; ok {
		info.subTopics = nil
		info.apis = nil
	}
	return nil
}

func (c *ControllerV1) OnMsg(agent com.AgentI, data []byte) ([]byte, error) {
	//解析opcode
	opCode, err := c.codec.ParseOpCode(data)
	if err != nil {
		return nil, err
	}
	//捕获一些事件
	switch opCode & com.MASK_ACTION {
	case com.ACTION_ROUTE: //路由
	case com.ACTION_PING: //ping
	}
	return data, nil
}
