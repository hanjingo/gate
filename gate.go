package gate

import (
	"errors"
	"sync"

	"github.com/hanjingo/gate/com"
	"github.com/hanjingo/network"
)

type Gate struct {
	id           uint8                           //gate id
	agents       map[interface{}]com.AgentI      //端点集合
	servers      map[interface{}]network.ServerI //tcp/ws/http服务器集合
	codecer      *codec                          //编解码器
	pluginSystem *PluginSystem                   //插件系统
}

func NewGate() *Gate {
	back := &Gate{
		agents:  make(map[interface{}]com.AgentI),
		servers: make(map[interface{}]network.ServerI),
		codecer: newCodec(),
	}
	return back
}

//跑起来
func (gate *Gate) Run(wg *sync.WaitGroup) {
	for _, s := range gate.servers {
		s.Run(wg)
	}
}

//处理新建连接
func (gate *Gate) onNewConn(c network.ConnI) {
	agent := newAgentV1(c)
	if gate.pluginSystem != nil {
		if err := gate.pluginSystem.onNewAgent(agent); err != nil {
			return
		}
	}
	gate.agents[agent.GetId()] = agent
}

//处理连接断开
func (gate *Gate) onConnClose(c network.ConnI) {
	agent, ok := gate.agents[c.GetId()]
	if !ok {
		//todo
		return
	}
	if gate.pluginSystem != nil {
		gate.pluginSystem.onAgentClose(agent)
	}
	//todo
}

//处理收到的消息
func (gate *Gate) handleMsg(agentId uint64, data []byte) error {
	agent, ok := gate.agents[agentId]
	if !ok {
		return errors.New("送信人不存在")
	}
	//编解码
	msg, err := gate.codecer.UnFormat(data)
	if err != nil || msg == nil {
		return err
	}
	return gate.pluginSystem.onMsg(agent, msg)
}
