package gate

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hanjingo/gate/com"
	"github.com/hanjingo/network"
	pc4 "github.com/hanjingo/protocol/v4"
)

type Gate struct {
	id           uint8                           //gate id
	agents       map[interface{}]com.AgentI      //端点集合
	servers      map[interface{}]network.ServerI //tcp/ws/http服务器集合
	pluginSystem *PluginSystem                   //插件系统
	codec        *pc4.Codec                      //编解码器
}

func NewGate(conf *GateConfig) *Gate {
	back := &Gate{
		agents:       make(map[interface{}]com.AgentI),
		servers:      make(map[interface{}]network.ServerI),
		pluginSystem: newPluginSystem(),
		codec:        pc4.NewCodec(),
	}
	for _, sconf := range conf.Servers {
		switch strings.ToUpper(sconf.Type) {
		case "WS":
			s, err := network.NewWsServer(sconf, back.onConnClose, back.onNewConn)
			if err != nil {
				panic(err)
			}
			back.servers[s.Name] = s
		}
	}
	return back
}

//跑起来
func (gate *Gate) Run(wg *sync.WaitGroup) {
	for id, s := range gate.servers {
		s.Run(wg)
		fmt.Println("服务器:", id, "已启动")
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
	delete(gate.agents, c.GetId())
}

//处理收到的消息
func (gate *Gate) handleMsg(agentId uint64, data []byte) error {
	agent, ok := gate.agents[agentId]
	if !ok {
		return errors.New("送信人不存在")
	}
	//解析操作码
	opcode, err := gate.codec.ParseOpCode(data)
	if err != nil {
		return err
	}
	//捕获系统消息
	switch opcode & com.MASK_ACTION {
	case com.ACTION_SYS:
		//todo
		return nil
	}
	return gate.pluginSystem.onMsg(agent, data)
}
