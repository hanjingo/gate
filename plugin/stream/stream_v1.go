package plugin

import (
	"errors"
	"time"

	"github.com/hanjingo/container"
	"github.com/hanjingo/gate/com"
)

type agentInfoV1 struct {
	id               interface{}     //id
	cache            *container.Pool //缓存
	life             time.Duration   //已使用的时长
	totalUsedStream  int64           //总共使用的流量(kb)
	usedStreamInDur  int64           //当前时间内已用流量(字节)
	totalStreamInDur int64           //当前时间内总可用流量
}
type StreamV1 struct {
	name        string
	agents      map[interface{}]*agentInfoV1
	streamSpeed int                  //流速(字节/s)
	scanDur     int                  //扫描间隔(s)
	reSendMsg   func(uint64, []byte) //重发消息函数
}

//streamSize:默认流速 	scanDur:扫描间隔	f:消息重发函数
func NewStreamV1(streamSpeed int, scanDur int, f func(uint64, []byte)) *StreamV1 {
	back := &StreamV1{
		name:        "StreamV1",
		agents:      make(map[interface{}]*agentInfoV1),
		streamSpeed: streamSpeed, //默认流速
		scanDur:     scanDur,     //扫描间隔
		reSendMsg:   f,
	}
	back.run()
	return back
}

func (s *StreamV1) Name() string {
	return s.name
}

func (s *StreamV1) OnNewAgent(agent com.AgentI) error {
	info := &agentInfoV1{
		id:               agent.GetId(),
		cache:            container.NewPool(-1),
		life:             0,
		totalUsedStream:  0,
		usedStreamInDur:  0,
		totalStreamInDur: int64(s.streamSpeed * s.scanDur),
	}
	s.agents[info.id] = info
	return nil
}

func (s *StreamV1) OnAgentClose(agent com.AgentI) error {
	delete(s.agents, agent.GetId()) //todo 暂时先这么做吧
	return nil
}

func (s *StreamV1) OnMsg(agent com.AgentI, data []byte) ([]byte, error) {
	if info, ok := s.agents[agent.GetId()]; ok {
		if info.totalStreamInDur-info.usedStreamInDur < int64(len(data)) {
			info.cache.Set(data)
			return nil, errors.New("流量已用完")
		}
	}
	return data, nil
}

//跑起来
func (s *StreamV1) run() {
	go func() {
		scanDur := time.Duration(s.scanDur) * time.Second
		tm := time.NewTimer(scanDur)
		for {
			select {
			case <-tm.C:
				for _, info := range s.agents {
					info.usedStreamInDur = 0
					info.life += scanDur
					//把消息取出来发给插件总线
					for data := info.cache.Get(); data != nil && s.reSendMsg != nil; {
						s.reSendMsg(info.id.(uint64), data.([]byte))
					}
				}
				tm.Reset(scanDur)
			}
		}
	}()
}
