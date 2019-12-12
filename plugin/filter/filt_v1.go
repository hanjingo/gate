package plugin

import "github.com/hanjingo/gate/com"

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
	return nil
}

func (f *FilterV1) OnAgentClose(agent com.AgentI) error {
	return nil
}

func (f *FilterV1) OnMsg(msg *com.Msg) (*com.Msg, error) {
	return msg, nil
}
