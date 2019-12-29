package filt_v1

import (
	"errors"

	"github.com/hanjingo/gate/com"
)

type agentInfo struct {
	id    interface{} //id
	bFilt bool        //是否被过滤
}
type Filt struct {
	name   string
	agents map[interface{}]*agentInfo
}

func NewFilt() *Filt {
	return &Filt{
		name:   "FiltV1",
		agents: make(map[interface{}]*agentInfo),
	}
}

func (f *Filt) onNewAgent(agent com.AgentI) error {
	info := &agentInfo{
		id:    agent.GetId(),
		bFilt: false,
	}
	f.agents[info.id] = info
	return nil
}

func (f *Filt) onAgentClose(agent com.AgentI) error {
	delete(f.agents, agent.GetId())
	return nil
}

func (f *Filt) onFilt(agent com.AgentI, data []byte) ([]byte, error) {
	if info, ok := f.agents[agent.GetId()]; ok {
		if info.bFilt {
			return nil, errors.New("消息已经被过滤")
		}
	}
	return data, nil
}
