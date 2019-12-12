package plugin

import "github.com/hanjingo/gate/com"

type agentInfoV1 struct {
	id   interface{} //id
	perm uint8       //权限
}
type PermV1 struct {
	name   string
	agents map[interface{}]*agentInfoV1
}

func NewPermV1() *PermV1 {
	return &PermV1{
		name:   "PermV1",
		agents: make(map[interface{}]*agentInfoV1),
	}
}

func (f *PermV1) Name() string {
	return f.name
}

func (f *PermV1) OnNewAgent(agent com.AgentI) error {
	return nil
}

func (f *PermV1) OnAgentClose(agent com.AgentI) error {
	return nil
}

func (f *PermV1) OnMsg(msg *com.Msg) (*com.Msg, error) {
	return msg, nil
}
