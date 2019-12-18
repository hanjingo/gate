package gate

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hanjingo/gate/com"
	"github.com/hanjingo/network"
	pv4 "github.com/hanjingo/protocol/v4"
)

const MASK_SYS = 0xa0000000 //系统权限

type Gate struct {
	id           uint8                           //gate id
	agents       map[interface{}]com.AgentI      //端点集合
	servers      map[interface{}]network.ServerI //tcp/ws/http服务器集合
	pluginSystem *PluginSystem                   //插件系统
	codec        *pv4.Codec                      //编解码器
}

func NewGate(conf *GateConfig) *Gate {
	back := &Gate{
		agents:       make(map[interface{}]com.AgentI),
		servers:      make(map[interface{}]network.ServerI),
		pluginSystem: newPluginSystem(),
		codec:        pv4.NewCodec(),
	}
	for _, sconf := range conf.Servers {
		switch strings.ToUpper(sconf.Type) {
		case "WS":
			s, err := network.NewWsServer(sconf, back.onConnClose, back.onNewConn, back.handleMsg)
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
func (gate *Gate) handleMsg(agentId uint64, data []byte) {
	agent, ok := gate.agents[agentId]
	if !ok {
		return
	}
	//解析操作码
	opcode, err := gate.codec.ParseOpCode(data)
	if err != nil {
		return
	}
	//捕获系统消息
	switch opcode & MASK_SYS {
	//todo
	}
	err = gate.pluginSystem.onMsg(agent, data)
}
