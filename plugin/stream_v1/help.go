package stream_v1

import (
	"github.com/hanjingo/gate/com"
	"github.com/hanjingo/golib/container"
)

func (s *StreamV1) reg() {
	if s.info == nil || s.info.CallBackMap == nil {
		return
	}
	m := s.info.CallBackMap
	m[com.AGENT_CONNECT] = s.onNewAgent
	m[com.AGENT_CLOSE] = s.onAgentClose
	m[OP_LIMIT_STREAM] = s.limit
	m[OP_SET_LIMIT_SPEED_CONF] = s.setSpeedConf
	m[OP_SET_LIMIT_SCANDUR_CONF] = s.setScanDurConf
}

func (s *StreamV1) onNewAgent(agent com.AgentI) {
	info := &agentInfo{
		id:               agent.GetId(),
		cache:            container.NewPool(-1),
		life:             0,
		totalUsedStream:  0,
		usedStreamInDur:  0,
		totalStreamInDur: int64(s.streamSpeed * s.scanDur),
	}
	s.agents()[info.id] = info
}

func (s *StreamV1) onAgentClose(agent com.AgentI) {
	delete(s.agents(), agent.GetId()) //todo 暂时先这么做吧
}

func (s *StreamV1) limit(agent com.AgentI, data []byte) {
	if info, ok := s.agents()[agent.GetId()]; ok {
		if info.totalStreamInDur-info.usedStreamInDur < int64(len(data)) {
			var temp []byte
			copy(temp, data)
			data = []byte{}
			info.cache.Set(temp)
			return
		}
	}
}

//设置流速
func (s *StreamV1) setSpeedConf(agent com.AgentI, n int) {
	s.streamSpeed = n
}

//设置扫描间隔
func (s *StreamV1) setScanDurConf(agent com.AgentI, n int) {
	s.scanDur = n
}
