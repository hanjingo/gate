package stream_v1

import (
	"errors"

	"github.com/hanjingo/container"
	"github.com/hanjingo/gate/com"
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

func (s *StreamV1) onNewAgent(agent com.AgentI) error {
	info := &agentInfo{
		id:               agent.GetId(),
		cache:            container.NewPool(-1),
		life:             0,
		totalUsedStream:  0,
		usedStreamInDur:  0,
		totalStreamInDur: int64(s.streamSpeed * s.scanDur),
	}
	s.agents()[info.id] = info
	return nil
}

func (s *StreamV1) onAgentClose(agent com.AgentI) error {
	delete(s.agents(), agent.GetId()) //todo 暂时先这么做吧
	return nil
}

func (s *StreamV1) limit(agent com.AgentI, data []byte) ([]byte, error) {
	if info, ok := s.agents()[agent.GetId()]; ok {
		if info.totalStreamInDur-info.usedStreamInDur < int64(len(data)) {
			info.cache.Set(data)
			return nil, errors.New("流量已用完")
		}
	}
	return data, nil
}

//设置流速
func (s *StreamV1) setSpeedConf(agent com.AgentI, n int) error {
	s.streamSpeed = n
	return nil
}

//设置扫描间隔
func (s *StreamV1) setScanDurConf(agent com.AgentI, n int) error {
	s.scanDur = n
	return nil
}
