package controller

import (
	"math/rand"
	"time"

	"github.com/hanjingo/util"

	"github.com/hanjingo/gate/com"
	pc4 "github.com/hanjingo/protocol/v4"
)

func (ctl *ControllerV1) regFunc() {
	ctl.actionFuncMap[com.ACTION_ROUTE] = ctl.onRoute //路由
	ctl.actionFuncMap[com.ACTION_PING] = ctl.onPing   //ping
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
	msg1 := &pc4.Messager{Content: &MsgPing{}}
	if err := ctl.codec.UnFormat(data, msg1); err != nil {
		return nil, err
	}
	//发送pong
	msg2 := &pc4.Messager{
		OpCode:    com.ACTION_PONG,
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

//处理请求信息
func (ctl *ControllerV1) onReq(agent com.AgentI, data []byte) ([]byte, error) {
	opcode, err := ctl.codec.ParseOpCode(data)
	if err != nil {
		return nil, err
	}
	recvs, err := ctl.codec.ParseRecv(data)
	if err != nil {
		return nil, err
	}
	if recvs != nil || len(recvs) > 0 {
		for _, recv := range recvs {
			if info, ok := ctl.agents[recv]; ok {
				info.agent.Write(data)
			}
		}
	} else { //如果不指定接受人就随机发一个
		if servs, ok := ctl.apiAgents[opcode]; ok {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			n := r.Intn(len(servs))
			for id, _ := range servs {
				n--
				if info, ok := ctl.agents[id]; ok && n < 1 {
					info.agent.Write(data)
				}
			}
		}
	}
	return data, nil
}

//处理订阅信息
func (ctl *ControllerV1) onSub(agent com.AgentI, data []byte) ([]byte, error) {
	msg := &pc4.Messager{Content: &MsgSub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		return data, err
	}
	content := msg.Content.(*MsgSub)
	if info, ok := ctl.agents[agent.GetId()]; ok {
		for _, topic := range content.Topics {
			info.subTopics[topic] = true
			if _, ok := ctl.subAgents[topic]; !ok {
				ctl.subAgents[topic] = make(map[interface{}]bool)
			}
			ctl.subAgents[topic][agent.GetId()] = true
		}
	}
	return data, nil
}

//处理取消订阅信息
func (ctl *ControllerV1) onUnSub(agent com.AgentI, data []byte) ([]byte, error) {
	msg := &pc4.Messager{Content: &MsgUnSub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		return data, err
	}
	content := msg.Content.(*MsgUnSub)
	if info, ok := ctl.agents[agent.GetId()]; ok {
		for _, topic := range content.Topics {
			delete(info.subTopics, topic)
			delete(ctl.subAgents[topic], agent.GetId())
		}
	}
	return data, nil
}

//处理发布消息
func (ctl *ControllerV1) onPub(agent com.AgentI, data []byte) ([]byte, error) {
	msg := &pc4.Messager{Content: &MsgPub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		return data, err
	}
	content := msg.Content.(*MsgPub)
	if recvs, ok := ctl.subAgents[content.Topic]; ok {
		for id, _ := range recvs {
			if info, ok := ctl.agents[id]; ok {
				info.agent.Write(data)
			}
		}
	}
	return data, nil
}
