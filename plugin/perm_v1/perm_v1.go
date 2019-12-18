package perm_v1

import (
	"errors"
	"strings"

	pv4 "github.com/hanjingo/protocol/v4"

	"github.com/hanjingo/gate/com"
)

const NORM_PERM = 0x20000000 //普通权限

type agentInfo struct {
	id   interface{} //id
	perm uint32      //权限
}
type Perm struct {
	name   string
	agents map[interface{}]*agentInfo
	codec  *pv4.Codec
}

func NewPerm() *Perm {
	return &Perm{
		name:   "PermV1",
		agents: make(map[interface{}]*agentInfo),
		codec:  pv4.NewCodec(),
	}
}

func (p *Perm) Name() string {
	return p.name
}

func (p *Perm) Set(name string, value interface{}) {
	switch strings.ToUpper(name) {
	case "NAME":
		if name, ok := value.(string); ok {
			p.name = name
		}
	}
}

func (p *Perm) OnNewAgent(agent com.AgentI) error {
	info := &agentInfo{
		id:   agent.GetId(),
		perm: NORM_PERM, //所有端点默认普通权限
	}
	p.agents[info.id] = info
	return nil
}

func (p *Perm) OnAgentClose(agent com.AgentI) error {
	delete(p.agents, agent.GetId())
	return nil
}

func (p *Perm) OnMsg(agent com.AgentI, data []byte) ([]byte, error) {
	//解析opcode
	opCode, err := p.codec.ParseOpCode(data)
	if err != nil {
		return nil, err
	}
	//检查权限
	if info, ok := p.agents[agent.GetId()]; ok {
		if opCode <= info.perm {
			return data, nil
		}
	}
	return nil, errors.New("权限不够")
}
