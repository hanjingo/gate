package controller

import (
	"time"

	"github.com/hanjingo/gate/com"
	pv4 "github.com/hanjingo/protocol/v4"
)

type agentInfoV1 struct {
	agent     com.AgentI
	id        interface{}     //id
	startTm   time.Time       //开始时间
	subTopics map[string]bool //订阅的主题
	apis      map[uint32]bool //可提供的api服务
}
type ControllerV1 struct {
	name          string
	agents        map[interface{}]*agentInfoV1
	actionFuncMap map[uint32]func(com.AgentI, []byte) ([]byte, error)
	apiAgents     map[uint32]map[interface{}]bool //提供api的端点集合
	subAgents     map[string]map[interface{}]bool //订阅了主题的端点集合 key:topic value:key:agentid value:
	codec         *pv4.Codec
}

func NewControllerV1() *ControllerV1 {
	back := &ControllerV1{
		name:          "ControllerV1",
		agents:        make(map[interface{}]*agentInfoV1),
		actionFuncMap: make(map[uint32]func(com.AgentI, []byte) ([]byte, error)),
		apiAgents:     make(map[uint32]map[interface{}]bool),
		subAgents:     make(map[string]map[interface{}]bool),
		codec:         pv4.NewCodec(),
	}
	back.regFunc()
	return back
}

func (c *ControllerV1) Name() string {
	return c.name
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
		tm := time.NewTimer(time.Duration(3000) * time.Millisecond)
		tempC := make(chan []byte, 1)
		defer close(tempC)
		select {
		case <-tm.C:
			return
		case tempC <- doRead(agent):
			data := <-tempC
			if data == nil {
				return
			}
			msg := &pv4.Messager{Content: &MsgNewAgentReq{}}
			if err := c.codec.UnFormat(data, msg); err != nil {
				return
			}
			rsp := &MsgNewAgentRsp{Result: true, Id: agent.GetId().(uint64)}
			msg1 := &pv4.Messager{
				OpCode:  com.ACTION_NEW_AGENT_RSP,
				Content: rsp,
			}
			data1, err := c.codec.Format(msg1)
			if err != nil {
				return
			}
			if err := agent.Write(data1); err != nil {
				info := &agentInfoV1{
					agent:     agent,
					id:        agent.GetId(),
					startTm:   time.Now(),
					subTopics: make(map[string]bool),
					apis:      make(map[uint32]bool),
				}
				c.agents[info.id] = info
			}
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
