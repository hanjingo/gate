package plugin

import (
	"errors"

	"github.com/hanjingo/gate/com"
)

type agentInfoV1 struct {
	id    interface{} //id
	bFilt bool        //是否被过滤
}
type FilterV1 struct {
	name   string
	agents map[interface{}]*agentInfoV1
}

func NewFilterV1() *FilterV1 {
	return &FilterV1{
		name:   "FilterV1",
		agents: make(map[interface{}]*agentInfoV1),
	}
}

func (f *FilterV1) Name() string {
	return f.name
}

func (f *FilterV1) OnNewAgent(agent com.AgentI) error {
	info := &agentInfoV1{
		id:    agent.GetId(),
		bFilt: false,
	}
	f.agents[info.id] = info
	return nil
}

func (f *FilterV1) OnAgentClose(agent com.AgentI) error {
	delete(f.agents, agent.GetId())
	return nil
}

func (f *FilterV1) OnMsg(agent com.AgentI, data []byte) ([]byte, error) {
	if info, ok := f.agents[agent.GetId()]; ok {
		if info.bFilt {
			return nil, errors.New("消息已经被过滤")
		}
	}
	return data, nil
}
