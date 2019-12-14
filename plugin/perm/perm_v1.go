package plugin

import (
	"errors"

	pc4 "github.com/hanjingo/protocol/v4"

	"github.com/hanjingo/gate/com"
)

type agentInfoV1 struct {
	id   interface{} //id
	perm uint32      //权限
}
type PermV1 struct {
	name   string
	agents map[interface{}]*agentInfoV1
	codec  *pc4.Codec
}

func NewPermV1() *PermV1 {
	return &PermV1{
		name:   "PermV1",
		agents: make(map[interface{}]*agentInfoV1),
		codec:  pc4.NewCodec(),
	}
}

func (p *PermV1) Name() string {
	return p.name
}

func (p *PermV1) OnNewAgent(agent com.AgentI) error {
	info := &agentInfoV1{
		id:   agent.GetId(),
		perm: com.PERM_NORMAL, //所有端点默认普通权限
	}
	p.agents[info.id] = info
	return nil
}

func (p *PermV1) OnAgentClose(agent com.AgentI) error {
	delete(p.agents, agent.GetId())
	return nil
}

func (p *PermV1) OnMsg(agent com.AgentI, data []byte) ([]byte, error) {
	//解析opcode
	opCode, err := p.codec.ParseOpCode(data)
	if err != nil {
		return nil, err
	}
	//检查权限
	if info, ok := p.agents[agent.GetId()]; ok {
		if com.MASK_ACTION&opCode <= info.perm {
			return data, nil
		}
	}
	return nil, errors.New("权限不够")
}
