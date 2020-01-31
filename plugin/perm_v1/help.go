package perm_v1

import (
	"errors"

	"github.com/hanjingo/gate/com"
)

func (p *PermV1) reg() {
	if p.info == nil || p.info.CallBackMap == nil {
		return
	}
	m := p.info.CallBackMap
	m[com.AGENT_CONNECT] = p.onNewAgent
	m[com.AGENT_CLOSE] = p.onAgentClose
	m[OP_CHANGE_PERM] = p.onCheckPerm
}

//端点建立
func (p *PermV1) onNewAgent(agent com.AgentI) error {
	info := &agentInfo{
		id:   agent.GetId(),
		perm: com.MASK_NORM, //所有端点默认普通权限
	}
	p.agents()[info.id] = info
	return nil
}

//端点断开
func (p *PermV1) onAgentClose(agent com.AgentI) error {
	delete(p.agents(), agent.GetId())
	return nil
}

//检查权限
func (p *PermV1) onCheckPerm(agent com.AgentI, data []byte) ([]byte, error) {
	//解析opcode
	opCode, err := p.codec.ParseOpCode(data)
	if err != nil {
		return nil, err
	}
	//检查权限
	if info, ok := p.agents()[agent.GetId()]; ok {
		if opCode <= info.perm {
			return data, nil
		}
	}
	return nil, errors.New("权限不够")
}
