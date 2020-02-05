package filt_v1

import (
	"github.com/hanjingo/gate/com"
)

func (f *FiltV1) reg() {
	if f.info == nil || f.info.CallBackMap == nil {
		return
	}
	m := f.info.CallBackMap
	m[com.AGENT_CONNECT] = f.onNewAgent    //建立端点
	m[com.AGENT_CLOSE] = f.onAgentClose    //端点断开
	m[OP_ON_FILT] = f.onFilt               //过滤
	m[OP_SET_FILT_AGENT] = f.onSetFilt     //设置过滤
	m[OP_SET_UNFILT_AGENT] = f.onSetUnFilt //设置不过滤
}

//建立端点
func (f *FiltV1) onNewAgent(agent com.AgentI) {
	info := &agentInfo{
		id:    agent.GetId(),
		bFilt: false,
	}
	f.agents()[info.id] = info
}

//端点断开
func (f *FiltV1) onAgentClose(agent com.AgentI) {
	delete(f.agents(), agent.GetId())
}

//过滤
func (f *FiltV1) onFilt(agent com.AgentI, data []byte) {
	if info, ok := f.agents()[agent.GetId()]; ok {
		if info.bFilt {
			data = []byte{}
			return
		}
	}
}

//设置过滤
func (f *FiltV1) onSetFilt(agent com.AgentI) {
	if info, ok := f.agents()[agent.GetId()]; ok {
		info.bFilt = true
	}
}

//解除过滤
func (f *FiltV1) onSetUnFilt(agent com.AgentI) {
	if info, ok := f.agents()[agent.GetId()]; ok {
		info.bFilt = false
	}
}
