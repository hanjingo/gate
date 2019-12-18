package control_v1

import (
	"strings"
	"time"

	"github.com/hanjingo/container"

	"github.com/hanjingo/gate/com"
	pv4 "github.com/hanjingo/protocol/v4"
)

type agentInfo struct {
	agent     com.AgentI
	id        interface{}    //id
	startTm   time.Time      //开始时间
	subTopics *container.Set //订阅的主题
	apis      *container.Set //可提供的api服务
}
type ControllerV1 struct {
	name                string
	agents              map[interface{}]*agentInfo
	actionFuncMap       map[uint32]func(com.AgentI, []byte) ([]byte, error)
	apiAgents           map[uint32]*container.Set //提供api的端点集合 key:api value:端点集合
	subAgents           map[string]*container.Set //订阅了主题的端点集合 key:topic value:端点集合
	codec               *pv4.Codec
	newAgentWaitTimeout int //新端点等待时长(ms)
}

func NewControllerV1() *ControllerV1 {
	back := &ControllerV1{
		name:                "ControllerV1",
		agents:              make(map[interface{}]*agentInfo),
		actionFuncMap:       make(map[uint32]func(com.AgentI, []byte) ([]byte, error)),
		apiAgents:           make(map[uint32]*container.Set),
		subAgents:           make(map[string]*container.Set),
		codec:               pv4.NewCodec(),
		newAgentWaitTimeout: 3000,
	}
	back.regFunc()
	return back
}

func (c *ControllerV1) Name() string {
	return c.name
}

func (s *ControllerV1) Set(name string, value interface{}) {
	switch strings.ToUpper(name) {
	case "NAME":
		if name, ok := value.(string); ok {
			s.name = name
		}
	}
}

func doRead(agent com.AgentI) []byte {
	data, err := agent.Read()
	if err != nil {
		return nil
	}
	return data
}

func (c *ControllerV1) OnNewAgent(agent com.AgentI) error {
	go func() {
		tm := time.NewTimer(time.Duration(c.newAgentWaitTimeout) * time.Millisecond)
		defer agent.Close()
		select {
		case <-tm.C:
			if !agent.IsValid() {
				agent.Close()
			}
			return
		}
	}()
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
	if f, ok := c.actionFuncMap[opCode]; ok {
		return f(agent, data)
	}
	return data, nil
}
