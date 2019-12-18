package control_v1

import (
	"time"

	"github.com/hanjingo/container"

	"github.com/hanjingo/util"

	"github.com/hanjingo/gate/com"
	pv4 "github.com/hanjingo/protocol/v4"
)

func (ctl *ControllerV1) regFunc() {
	ctl.actionFuncMap[OP_NEW_AGENT] = ctl.onNewAgentReq //连接验证
	ctl.actionFuncMap[OP_ROUTE] = ctl.onRoute           //路由
	ctl.actionFuncMap[OP_PING] = ctl.onPing             //ping
	ctl.actionFuncMap[OP_SUB] = ctl.onSub               //订阅
	ctl.actionFuncMap[OP_UNSUB] = ctl.onUnSub           //取消订阅
	ctl.actionFuncMap[OP_PUB] = ctl.onPub               //发布
}

//连接验证
func (ctl *ControllerV1) onNewAgentReq(agent com.AgentI, data []byte) ([]byte, error) {
	msg := &pv4.Messager{Content: &MsgNewAgentReq{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		return nil, err
	}
	//验证 todo

	rsp := &MsgNewAgentRsp{Result: true, Id: agent.GetId().(uint64)}
	msg1 := &pv4.Messager{
		OpCode:  OP_NEW_AGENT_RSP,
		Content: rsp,
	}
	data1, err := ctl.codec.Format(msg1)
	if err != nil {
		return nil, err
	}
	if err := agent.Write(data1); err != nil {
		info := &agentInfo{
			agent:     agent,
			id:        agent.GetId(),
			startTm:   time.Now(),
			subTopics: container.NewSet(),
			apis:      container.NewSet(),
		}
		ctl.agents[info.id] = info
	}
	return nil, nil
}

//处理路由
func (ctl *ControllerV1) onRoute(agent com.AgentI, data []byte) ([]byte, error) {
	recvs, err := ctl.codec.ParseRecv(data)
	if err != nil {
		return nil, err
	}
	for _, recv := range recvs {
		if info, ok := ctl.agents[recv]; ok {
			info.agent.Write(data)
		}
	}
	return data, nil
}

//处理ping
func (ctl *ControllerV1) onPing(agent com.AgentI, data []byte) ([]byte, error) {
	//解析ping
	msg1 := &pv4.Messager{Content: &MsgPing{}}
	if err := ctl.codec.UnFormat(data, msg1); err != nil {
		return nil, err
	}
	//发送pong
	msg2 := &pv4.Messager{
		OpCode:    OP_PONG,
		Receiver:  []uint64{msg1.Sender},
		TimeStamp: util.TimeDurToMilliSecond(time.Now().Sub(com.START_TIME)),
		Content:   &MsgPong{},
	}
	data2, err := ctl.codec.Format(msg2)
	if err != nil {
		return nil, err
	}
	agent.Write(data2)
	return data, nil
}

//处理订阅信息
func (ctl *ControllerV1) onSub(agent com.AgentI, data []byte) ([]byte, error) {
	msg := &pv4.Messager{Content: &MsgSub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		return data, err
	}
	content := msg.Content.(*MsgSub)
	if _, ok := ctl.agents[agent.GetId()]; ok {
		for _, topic := range content.Topics {
			if _, ok := ctl.subAgents[topic]; !ok {
				ctl.subAgents[topic] = container.NewSet()
			}
			ctl.subAgents[topic].Add(agent.GetId())
		}
	}
	return data, nil
}

//处理取消订阅信息
func (ctl *ControllerV1) onUnSub(agent com.AgentI, data []byte) ([]byte, error) {
	msg := &pv4.Messager{Content: &MsgUnSub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		return data, err
	}
	content := msg.Content.(*MsgUnSub)
	if info, ok := ctl.agents[agent.GetId()]; ok {
		for _, topic := range content.Topics {
			info.subTopics.Del(topic)
			if agents, ok := ctl.subAgents[topic]; ok {
				agents.Del(agent.GetId())
				if agents.Len() == 0 {
					delete(ctl.subAgents, topic)
				}
			}
		}
	}
	return data, nil
}

//处理发布消息
func (ctl *ControllerV1) onPub(agent com.AgentI, data []byte) ([]byte, error) {
	msg := &pv4.Messager{Content: &MsgPub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		return data, err
	}
	content := msg.Content.(*MsgPub)
	if recvs, ok := ctl.subAgents[content.Topic]; ok {
		recvs.Range(func(recv interface{}) {
			info, ok := recv.(*agentInfo)
			if !ok {
				return
			}
			info.agent.Write(data)
		})
	}
	return data, nil
}
