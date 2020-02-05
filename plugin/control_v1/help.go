package control_v1

import (
	"time"

	"github.com/hanjingo/container"

	"github.com/hanjingo/gate/com"
	pv4 "github.com/hanjingo/protocol/v4"
)

func (ctl *ControllerV1) reg() {
	if ctl.info == nil || ctl.info.CallBackMap == nil {
		return
	}
	m := ctl.info.CallBackMap
	m[com.AGENT_CONNECT] = ctl.onAgentConnect //端点连接
	m[com.AGENT_CLOSE] = ctl.onAgentClose     //关闭端点
	m[OP_NEW_AGENT_REQ] = ctl.onNewAgentReq   //建立端点请求
	m[OP_ROUTE] = ctl.onRoute                 //处理路由
	m[OP_PING] = ctl.onPing                   //处理ping
	m[OP_SUB] = ctl.onSub                     //处理订阅
	m[OP_UNSUB] = ctl.onUnSub                 //处理取消订阅
	m[OP_PUB] = ctl.onPub                     //处理发布
}

//预处理
func callBefore(agent com.AgentI, data []byte) {
	//todo
}

//新建agent
func (c *ControllerV1) onAgentConnect(agent com.AgentI) {
	log.Debug("onAgentConnect>> opcode:%b", com.AGENT_CONNECT)
	go func() {
		tm := time.NewTimer(time.Duration(c.newAgentWaitTimeout) * time.Millisecond)
		tm1 := time.NewTimer(time.Duration(3000 * time.Millisecond))
		defer agent.Close()
		select {
		case <-tm.C:
			if !agent.IsValid() {
				agent.Close()
			}
			return
		case <-tm1.C:
			return
		}
	}()

	ntf := &MsgAgentConnSucc{Id: agent.GetId().(uint64)}
	msg := &pv4.Messager{
		OpCode:  OP_AGENT_CONNECT_SUCCESS,
		Content: ntf,
	}
	data, err := c.codec.Format(msg)
	if err != nil {
		log.Error("连接建立时,格式化返回消息失败,错误:%v", err)
		return
	}
	if err := agent.Write(data); err != nil {
		log.Error("连接建立时,返回消息失败,错误:%v", err)
	}
}

//关闭agent
func (c *ControllerV1) onAgentClose(agent com.AgentI) {
	log.Debug("onAgentClose>>")
	if info, ok := c.agents()[agent.GetId()]; ok {
		info.subTopics = nil
		info.apis = nil
	}
}

//连接验证
func (ctl *ControllerV1) onNewAgentReq(agent com.AgentI, data []byte) {
	log.Debug("onNewAgentReq>>")
	msg := &pv4.Messager{Content: &MsgNewAgentReq{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		log.Error("连接验证时,反格式化失败,错误:%v", err)
		return
	}
	//验证 todo
	rsp := &MsgNewAgentRsp{Result: true, Id: agent.GetId().(uint64)}
	msg1 := &pv4.Messager{
		OpCode:  OP_NEW_AGENT_RSP,
		Content: rsp,
	}
	data1, err := ctl.codec.Format(msg1)
	if err != nil {
		log.Error("连接验证时,格式化返回失败,错误:%v", err)
		return
	}
	if err := agent.Write(data1); err != nil {
		info := &agentInfo{
			agent:     agent,
			id:        agent.GetId(),
			startTm:   time.Now(),
			subTopics: container.NewSet(),
			apis:      container.NewSet(),
		}
		ctl.agents()[info.id] = info
	}
}

//处理路由
func (ctl *ControllerV1) onRoute(agent com.AgentI, data []byte) {
	log.Debug("onRoute>>")
	recvs, err := ctl.codec.ParseRecv(data)
	if err != nil {
		log.Debug("路由时解析请求失败,错误:%v", err)
		return
	}
	for _, recv := range recvs {
		if info, ok := ctl.agents()[recv]; ok {
			info.agent.Write(data)
		}
	}
}

//处理ping
func (ctl *ControllerV1) onPing(agent com.AgentI, data []byte) {
	log.Debug("onPing>>")
	//解析ping
	msg1 := &pv4.Messager{Content: &MsgPing{}}
	if err := ctl.codec.UnFormat(data, msg1); err != nil {
		log.Error("处理ping时反格式化失败,错误:%v", err)
		return
	}
	//发送pong
	msg2 := &pv4.Messager{
		OpCode:   OP_PONG,
		Receiver: []uint64{msg1.Sender},
		Content:  &MsgPong{},
	}
	data2, err := ctl.codec.Format(msg2)
	if err != nil {
		log.Error("处理ping时格式化返回失败,错误:%v", err)
		return
	}
	agent.Write(data2)
}

//处理订阅信息
func (ctl *ControllerV1) onSub(agent com.AgentI, data []byte) {
	log.Debug("onSub>>")
	msg := &pv4.Messager{Content: &MsgSub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		log.Error("处理订阅时,反格式化请求失败,错误:%v", err)
		return
	}
	content := msg.Content.(*MsgSub)
	if _, ok := ctl.agents()[agent.GetId()]; ok {
		for _, topic := range content.Topics {
			if _, ok := ctl.subAgents[topic]; !ok {
				ctl.subAgents[topic] = container.NewSet()
			}
			ctl.subAgents[topic].Add(agent.GetId())
		}
	}
}

//处理取消订阅信息
func (ctl *ControllerV1) onUnSub(agent com.AgentI, data []byte) {
	log.Debug("onUnSub>>")
	msg := &pv4.Messager{Content: &MsgUnSub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		log.Error("取消订阅时反格式化失败,错误:%v", err)
		return
	}
	content := msg.Content.(*MsgUnSub)
	if info, ok := ctl.agents()[agent.GetId()]; ok {
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
}

//处理发布消息
func (ctl *ControllerV1) onPub(agent com.AgentI, data []byte) {
	log.Debug("onPub>>")
	msg := &pv4.Messager{Content: &MsgPub{}}
	if err := ctl.codec.UnFormat(data, msg); err != nil {
		log.Error("处理发布时解析请求失败,错误:%v", err)
		return
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
}
